# bash completion for go                                   -*- shell-script -*-

_new_split_longopt()
{
    if [[ "$cur" == -?*=* ]]; then
        # Cut also backslash before '=' in case it ended up there
        # for some reason.
        prev="${cur%%?(\\)=*}"
        cur="${cur#*=}"
        return 0
    fi

    return 1
}

_make_build_flags()
{
    if [[ $cword -gt 1 ]]; then
        case $prev in
            -buildmode)
                comms="archive c-archive c-shared default shared exe pie plugin"
                COMPREPLY=($(compgen -W "$comms" -- "$cur"))
                ;;
            -compiler)
                comms="gc gccgo"
                COMPREPLY=($(compgen -W "$comms" -- "$cur"))
                ;;
            -pkgdir)
                _filedir -d
                ;;
            -p|-asmflags|-gccgoflags|-gcflags|-installsuffix|-ldflags|-mod|-tags|-toolexec)
                ;;
            *)
                return 1
                ;;
        esac
    fi
    return 0
}

_parse_symbols()
{
    local sym=$1
    if [[ $1 == *\.* ]]; then
        sym=${1%\.*}
    else
        sym=""
    fi
    [[ ${sym-} ]] && sym_prefix="${sym}."

    [[ ${sym##*\.} == *[A-Z]* ]] || all="-all"

    local block_open=false
    local type_open=false

    go doc $all $sym 2>&1 | while read -r line; do

        if [[ $block_open == true ]]; then
            if [[ $line =~ ^\)$ ]]; then
                block_open=false
                continue
            elif [[ ${all-} ]] && [[ $line =~ ^[[:space:]]?([A-Z][0-9A-Za-z_]+)([[:space:]]|$) ]]; then
                match_result="${BASH_REMATCH[1]}"
            else
                continue
            fi
        elif [[ $type_open == true ]]; then
            if [[ $line =~ ^\}$ ]]; then
                type_open=false
                continue
            elif [[ $line =~ ^[[:space:]]?([A-Z][0-9A-Za-z_]+)([[:space:]]|\() ]]; then
                match_result="${BASH_REMATCH[1]}"
            else
                continue
            fi
        elif [[ $line =~ (type|const|var)[[:space:]]+([A-Z][0-9A-Za-z_]+)[[:space:]]+ ]]; then            
            [[ ${all-} ]] || if [[ $line =~ type.*\{ ]]; then
                type_open=true
                continue
            fi

            [[ ${all-} ]] && match_result="${BASH_REMATCH[2]}"
        elif [[ $line =~ \
            func([[:space:]]+\(.*\)[[:space:]]+|[[:space:]]+)([A-Z][0-9A-Za-z]+)[[:space:]]?\( ]]; then
            [[ ${all-} ]] && match_result="${BASH_REMATCH[2]}"
        elif [[ $line =~ (var|const)[[:space:]]+\( ]]; then
            block_open=true
            continue
        else
            continue
        fi
        
        if [[ $match_result != ${sym##*\.} ]]; then
            printf "%s\n" "${sym_prefix}${match_result}"
        fi
    done
}

_go()
{
    local cur prev words cword split
    _init_completion -s || return

    _new_split_longopt && split=true

    basic_comms="bug build clean doc env fix fmt generate get help install list mod run test tool version vet"
    build_flags=$(_parse_help go "help build")
    if [[ $cword == 1 ]]; then
        comms=$basic_comms
        COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
    else
        case ${words[1]} in
            bug)
                ;;
            build)
                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="${build_flags} -o -i"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                fi

                _make_build_flags && return

                go_pkgs=$(go list ./...)
                pkg_names=( $(compgen -W "$go_pkgs" -- "$cur") )

                case $prev in
                    -o)
                        _filedir
                        ;;
                    *)
                        _filedir "@(go)"
                        ;;
                esac

                compopt +o filenames 2>/dev/null
                COMPREPLY+=( ${pkg_names[@]} )
                ;;    
            clean)
                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="${build_flags} -i -n -r -x -cache -testcache-modcache"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                fi

                _make_build_flags && return

                go_pkgs=$(go list ./...)
                COMPREPLY+=( $(compgen -W "$go_pkgs" -- "$cur") )
                ;;
            doc)
                if [[ $cword -gt 2 ]] && [[ $prev == [a-z]* ]]; then
                    if [[ $cur == *\.* ]]; then
                        pref="${prev}.${cur}"
                    else
                        pref="${prev}."
                    fi
                    syms=$(_parse_symbols "$pref")
                    comms=${syms//"${prev}."/}
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                    return
                fi

                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="-u -c -all -cmd -src"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                    return
                fi

                comms=$(_parse_symbols "$cur")
                pkg_paths=". $(go env GOROOT) $(go env GOPATH)"
                go_pkgs=$(go list -f "{{.Name}} {{.ImportPath}}" ... 2>/dev/null)
                go_pkgs=${go_pkgs//main /}
                go list 2>/dev/null && go_pkgs="${go_pkgs} $(go list -f '{{.Name}} {{.ImportPath}}' ./... 2>/dev/null)"
                COMPREPLY=( $(compgen -W "$comms $go_pkgs" -- "$cur") )
                ;;
            env)
                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="-json -u -w"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                    return
                fi

                write_flag=false
                for i in "${words[@]}"; do
                    [[ $i == "-w" ]] && write_flag=true
                done

                if [[ $write_flag == true ]]; then

                    if [[ $cur == *=* ]]; then
                        cur=${cur#*=}
                        _filedir
                    else
                        comms=$(go env 2>&1 | while read -r line; do \
                            printf "%s\n" "${line%%=*}="; \
                        done)
                        COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                        [[ $COMPREPLY == *= ]] && compopt -o nospace
                    fi
                else
                    comms=$(go env 2>&1 | while read -r line; do \
                        printf "%s\n" "${line%%=*}"; \
                    done)
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                fi
                ;;
            fix)
                go_pkgs=$(go list ./...)
                COMPREPLY=( $(compgen -W "$go_pkgs" -- "$cur") )
                ;;
            fmt)
                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="-n -x"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                    return
                fi

                go_pkgs=$(go list ./...)
                COMPREPLY=( $(compgen -W "$go_pkgs" -- "$cur") )
                ;;
            generate)
                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="${build_flags} -run -n -v -x"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                    return
                fi
                case $prev in
                    -run)
                        ;;
                    *)
                        go_pkgs=$(go list ./...)
                        pkg_names=( $(compgen -W "$go_pkgs" -- "$cur") )
                        _filedir "@(go)"
                        compopt +o filenames 2>/dev/null
                        COMPREPLY+=( ${pkg_names[@]} )
                        ;;
                esac
                ;;
            get)
                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="${build_flags} -d -t -u -v -insecure"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                    return
                fi
                
                go_pkgs=$(go list ./...)
                COMPREPLY=( $(compgen -W "$go_pkgs" -- "$cur") )
                ;;
            help)
                comms="${basic_comms/help/} buildmode c cache environment filetype go.mod gopath gopath-get goproxy importpath modules module-get module-auth module-private packages testflag testfunc"
                COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                ;;
            install)
                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="${build_flags} -i"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                    return
                fi

                go_pkgs=$(go list ./...)
                COMPREPLY=( $(compgen -W "$go_pkgs" -- "$cur") )
                ;;
            list)
                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="${build_flags} -f -json -m -compiled -deps -e -export -find -test -u"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                    return
                fi

                case $prev in
                    -f)
                        ;;
                    *)
                        go_pkgs=$(go list ./...)
                        pkg_names=( $(compgen -W "$go_pkgs" -- "$cur") )
                        _filedir "@(go)"
                        compopt +o filenames 2>/dev/null
                        COMPREPLY+=( ${pkg_names[@]} )
                        ;;
                esac
                ;;
            mod)
                comms="download edit graph init tidy vendor verify why"
                COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                ;;
            run)
                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="${build_flags} -exec"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                    return
                fi

                case $prev in
                    -exec)
                        ;;
                    *)
                        go_pkgs="$(go list ./...)"
                        COMPREPLY=( $(compgen -W "$go_pkgs" -- "$cur") )
                        ;;
                esac
                ;;
            test)
                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="-args -c -exec -i -json -o"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                    return
                fi

                go_pkgs=$(go list ./...)
                pkg_names=( $(compgen -W "$go_pkgs" -- "$cur") )

                case $prev in
                    -o)
                        _filedir
                        compopt +o filenames 2>/dev/null
                        COMPREPLY+=( ${pkg_names[@]} )
                        ;;
                    -args|-exec)
                        ;;
                    *)
                        _filedir "@(go)"
                        compopt +o filenames 2>/dev/null
                        COMPREPLY+=( ${pkg_names[@]} )
                        ;;
                esac
                ;;
            tool)
                ;;
            version)
                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="-m -v"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                    return
                fi

                _filedir "@(go)"
                ;;
            vet)
                if [[ -n $cur ]] && [[ $cur == -* ]]; then
                    comms="${build_flags} -n -x -vettool"
                    COMPREPLY=( $(compgen -W "$comms" -- "$cur") )
                    return
                fi

                case $prev in
                    -vettool)
                        ;;
                    *)
                        go_pkgs=$(go list ./...)
                        COMPREPLY=( $(compgen -W "$go_pkgs" -- "$cur") )
                        ;;
                esac
                ;;
            *)
                ;;
        esac
    fi
    return
} &&
complete -F _go go

# ex: filetype=sh