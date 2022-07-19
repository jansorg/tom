# bash completion for tom                                  -*- shell-script -*-

__tom_debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE:-} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__tom_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__tom_index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__tom_contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__tom_handle_go_custom_completion()
{
    __tom_debug "${FUNCNAME[0]}: cur is ${cur}, words[*] is ${words[*]}, #words[@] is ${#words[@]}"

    local shellCompDirectiveError=1
    local shellCompDirectiveNoSpace=2
    local shellCompDirectiveNoFileComp=4
    local shellCompDirectiveFilterFileExt=8
    local shellCompDirectiveFilterDirs=16

    local out requestComp lastParam lastChar comp directive args

    # Prepare the command to request completions for the program.
    # Calling ${words[0]} instead of directly tom allows to handle aliases
    args=("${words[@]:1}")
    # Disable ActiveHelp which is not supported for bash completion v1
    requestComp="TOM_ACTIVE_HELP=0 ${words[0]} __completeNoDesc ${args[*]}"

    lastParam=${words[$((${#words[@]}-1))]}
    lastChar=${lastParam:$((${#lastParam}-1)):1}
    __tom_debug "${FUNCNAME[0]}: lastParam ${lastParam}, lastChar ${lastChar}"

    if [ -z "${cur}" ] && [ "${lastChar}" != "=" ]; then
        # If the last parameter is complete (there is a space following it)
        # We add an extra empty parameter so we can indicate this to the go method.
        __tom_debug "${FUNCNAME[0]}: Adding extra empty parameter"
        requestComp="${requestComp} \"\""
    fi

    __tom_debug "${FUNCNAME[0]}: calling ${requestComp}"
    # Use eval to handle any environment variables and such
    out=$(eval "${requestComp}" 2>/dev/null)

    # Extract the directive integer at the very end of the output following a colon (:)
    directive=${out##*:}
    # Remove the directive
    out=${out%:*}
    if [ "${directive}" = "${out}" ]; then
        # There is not directive specified
        directive=0
    fi
    __tom_debug "${FUNCNAME[0]}: the completion directive is: ${directive}"
    __tom_debug "${FUNCNAME[0]}: the completions are: ${out}"

    if [ $((directive & shellCompDirectiveError)) -ne 0 ]; then
        # Error code.  No completion.
        __tom_debug "${FUNCNAME[0]}: received error from custom completion go code"
        return
    else
        if [ $((directive & shellCompDirectiveNoSpace)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __tom_debug "${FUNCNAME[0]}: activating no space"
                compopt -o nospace
            fi
        fi
        if [ $((directive & shellCompDirectiveNoFileComp)) -ne 0 ]; then
            if [[ $(type -t compopt) = "builtin" ]]; then
                __tom_debug "${FUNCNAME[0]}: activating no file completion"
                compopt +o default
            fi
        fi
    fi

    if [ $((directive & shellCompDirectiveFilterFileExt)) -ne 0 ]; then
        # File extension filtering
        local fullFilter filter filteringCmd
        # Do not use quotes around the $out variable or else newline
        # characters will be kept.
        for filter in ${out}; do
            fullFilter+="$filter|"
        done

        filteringCmd="_filedir $fullFilter"
        __tom_debug "File filtering command: $filteringCmd"
        $filteringCmd
    elif [ $((directive & shellCompDirectiveFilterDirs)) -ne 0 ]; then
        # File completion for directories only
        local subdir
        # Use printf to strip any trailing newline
        subdir=$(printf "%s" "${out}")
        if [ -n "$subdir" ]; then
            __tom_debug "Listing directories in $subdir"
            __tom_handle_subdirs_in_dir_flag "$subdir"
        else
            __tom_debug "Listing directories in ."
            _filedir -d
        fi
    else
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${out}" -- "$cur")
    fi
}

__tom_handle_reply()
{
    __tom_debug "${FUNCNAME[0]}"
    local comp
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            while IFS='' read -r comp; do
                COMPREPLY+=("$comp")
            done < <(compgen -W "${allflags[*]}" -- "$cur")
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%=*}"
                __tom_index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION:-}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi

            if [[ -z "${flag_parsing_disabled}" ]]; then
                # If flag parsing is enabled, we have completed the flags and can return.
                # If flag parsing is disabled, we may not know all (or any) of the flags, so we fallthrough
                # to possibly call handle_go_custom_completion.
                return 0;
            fi
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __tom_index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions+=("${must_have_one_noun[@]}")
    elif [[ -n "${has_completion_function}" ]]; then
        # if a go completion function is provided, defer to that function
        __tom_handle_go_custom_completion
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    while IFS='' read -r comp; do
        COMPREPLY+=("$comp")
    done < <(compgen -W "${completions[*]}" -- "$cur")

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        while IFS='' read -r comp; do
            COMPREPLY+=("$comp")
        done < <(compgen -W "${noun_aliases[*]}" -- "$cur")
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
        if declare -F __tom_custom_func >/dev/null; then
            # try command name qualified custom func
            __tom_custom_func
        else
            # otherwise fall back to unqualified for compatibility
            declare -F __custom_func >/dev/null && __custom_func
        fi
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi

    # If there is only 1 completion and it is a flag with an = it will be completed
    # but we don't want a space after the =
    if [[ "${#COMPREPLY[@]}" -eq "1" ]] && [[ $(type -t compopt) = "builtin" ]] && [[ "${COMPREPLY[0]}" == --*= ]]; then
       compopt -o nospace
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__tom_handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__tom_handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1 || return
}

__tom_handle_flag()
{
    __tom_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue=""
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __tom_debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __tom_contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __tom_contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    # flaghash variable is an associative array which is only supported in bash > 3.
    if [[ -z "${BASH_VERSION:-}" || "${BASH_VERSINFO[0]:-}" -gt 3 ]]; then
        if [ -n "${flagvalue}" ] ; then
            flaghash[${flagname}]=${flagvalue}
        elif [ -n "${words[ $((c+1)) ]}" ] ; then
            flaghash[${flagname}]=${words[ $((c+1)) ]}
        else
            flaghash[${flagname}]="true" # pad "true" for bool flag
        fi
    fi

    # skip the argument to a two word flag
    if [[ ${words[c]} != *"="* ]] && __tom_contains_word "${words[c]}" "${two_word_flags[@]}"; then
        __tom_debug "${FUNCNAME[0]}: found a flag ${words[c]}, skip the next argument"
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__tom_handle_noun()
{
    __tom_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __tom_contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __tom_contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__tom_handle_command()
{
    __tom_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_tom_root_command"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __tom_debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__tom_handle_word()
{
    if [[ $c -ge $cword ]]; then
        __tom_handle_reply
        return
    fi
    __tom_debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __tom_handle_flag
    elif __tom_contains_word "${words[c]}" "${commands[@]}"; then
        __tom_handle_command
    elif [[ $c -eq 0 ]]; then
        __tom_handle_command
    elif __tom_contains_word "${words[c]}" "${command_aliases[@]}"; then
        # aliashash variable is an associative array which is only supported in bash > 3.
        if [[ -z "${BASH_VERSION:-}" || "${BASH_VERSINFO[0]:-}" -gt 3 ]]; then
            words[c]=${aliashash[${words[c]}]}
            __tom_handle_command
        else
            __tom_handle_noun
        fi
    else
        __tom_handle_noun
    fi
    __tom_handle_word
}

__tom_projects_get()
{
    local -a projects
    readarray -t COMPREPLY < <(tom projects 2>/dev/null | grep "$cur" | sed -e 's/ /\\ /g')
}

__tom_get_projects()
{
    if [[ ${#nouns[@]} -eq 0 ]]; then
        __tom_projects_get ""
	else
	    __tom_projects_get ${nouns[${#nouns[@]} -1]}
    fi
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__custom_func() {
    case ${last_command} in
        tom_start)
            __tom_get_projects
            return
            ;;
        *)
            ;;
    esac
}

_tom_cancel()
{
    last_command="tom_cancel"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_config_set()
{
    last_command="tom_config_set"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    must_have_one_noun+=("activity.stop_on_start")
    must_have_one_noun+=("backup.directory")
    must_have_one_noun+=("backup.max_to_keep")
    must_have_one_noun+=("data_dir")
    must_have_one_noun+=("projects.create_missing")
    noun_aliases=()
}

_tom_config()
{
    last_command="tom_config"

    command_aliases=()

    commands=()
    commands+=("set")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--output=")
    two_word_flags+=("--output")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output")
    local_nonpersistent_flags+=("--output=")
    local_nonpersistent_flags+=("-o")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_create_project()
{
    last_command="tom_create_project"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--output=")
    two_word_flags+=("--output")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output")
    local_nonpersistent_flags+=("--output=")
    local_nonpersistent_flags+=("-o")
    flags+=("--parent=")
    two_word_flags+=("--parent")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--parent")
    local_nonpersistent_flags+=("--parent=")
    local_nonpersistent_flags+=("-p")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_create_tag()
{
    last_command="tom_create_tag"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_create()
{
    last_command="tom_create"

    command_aliases=()

    commands=()
    commands+=("project")
    commands+=("tag")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_edit_frame()
{
    last_command="tom_edit_frame"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--archived")
    local_nonpersistent_flags+=("--archived")
    flags+=("--end=")
    two_word_flags+=("--end")
    two_word_flags+=("-t")
    local_nonpersistent_flags+=("--end")
    local_nonpersistent_flags+=("--end=")
    local_nonpersistent_flags+=("-t")
    flags+=("--name-delimiter=")
    two_word_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--notes=")
    two_word_flags+=("--notes")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--notes")
    local_nonpersistent_flags+=("--notes=")
    local_nonpersistent_flags+=("-n")
    flags+=("--project=")
    two_word_flags+=("--project")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--project")
    local_nonpersistent_flags+=("--project=")
    local_nonpersistent_flags+=("-p")
    flags+=("--start=")
    two_word_flags+=("--start")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--start")
    local_nonpersistent_flags+=("--start=")
    local_nonpersistent_flags+=("-f")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_edit_project()
{
    last_command="tom_edit_project"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--hourly-rate=")
    two_word_flags+=("--hourly-rate")
    local_nonpersistent_flags+=("--hourly-rate")
    local_nonpersistent_flags+=("--hourly-rate=")
    flags+=("--name=")
    two_word_flags+=("--name")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--name")
    local_nonpersistent_flags+=("--name=")
    local_nonpersistent_flags+=("-n")
    flags+=("--name-delimiter=")
    two_word_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--note-required=")
    two_word_flags+=("--note-required")
    local_nonpersistent_flags+=("--note-required")
    local_nonpersistent_flags+=("--note-required=")
    flags+=("--parent=")
    two_word_flags+=("--parent")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--parent")
    local_nonpersistent_flags+=("--parent=")
    local_nonpersistent_flags+=("-p")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_edit()
{
    last_command="tom_edit"

    command_aliases=()

    commands=()
    commands+=("frame")
    commands+=("project")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_frames_archive()
{
    last_command="tom_frames_archive"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--include-subprojects")
    local_nonpersistent_flags+=("--include-subprojects")
    flags+=("--name-delimiter=")
    two_word_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--project=")
    two_word_flags+=("--project")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--project")
    local_nonpersistent_flags+=("--project=")
    local_nonpersistent_flags+=("-p")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_frames()
{
    last_command="tom_frames"

    command_aliases=()

    commands=()
    commands+=("archive")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--archived")
    local_nonpersistent_flags+=("--archived")
    flags+=("--delimiter=")
    two_word_flags+=("--delimiter")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--delimiter")
    local_nonpersistent_flags+=("--delimiter=")
    local_nonpersistent_flags+=("-d")
    flags+=("--format=")
    two_word_flags+=("--format")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    local_nonpersistent_flags+=("-f")
    flags+=("--output=")
    two_word_flags+=("--output")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output")
    local_nonpersistent_flags+=("--output=")
    local_nonpersistent_flags+=("-o")
    flags+=("--project=")
    two_word_flags+=("--project")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--project")
    local_nonpersistent_flags+=("--project=")
    local_nonpersistent_flags+=("-p")
    flags+=("--subprojects")
    flags+=("-s")
    local_nonpersistent_flags+=("--subprojects")
    local_nonpersistent_flags+=("-s")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_help()
{
    last_command="tom_help"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    has_completion_function=1
    noun_aliases=()
}

_tom_import_fanurio()
{
    last_command="tom_import_fanurio"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_import_macTimeTracker()
{
    last_command="tom_import_macTimeTracker"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_import_watson()
{
    last_command="tom_import_watson"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_import()
{
    last_command="tom_import"

    command_aliases=()

    commands=()
    commands+=("fanurio")
    commands+=("macTimeTracker")
    commands+=("watson")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_projects()
{
    last_command="tom_projects"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--delimiter=")
    two_word_flags+=("--delimiter")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--delimiter")
    local_nonpersistent_flags+=("--delimiter=")
    local_nonpersistent_flags+=("-d")
    flags+=("--format=")
    two_word_flags+=("--format")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    local_nonpersistent_flags+=("-f")
    flags+=("--name-delimiter=")
    two_word_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--output=")
    two_word_flags+=("--output")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output")
    local_nonpersistent_flags+=("--output=")
    local_nonpersistent_flags+=("-o")
    flags+=("--recent=")
    two_word_flags+=("--recent")
    local_nonpersistent_flags+=("--recent")
    local_nonpersistent_flags+=("--recent=")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_remove_all()
{
    last_command="tom_remove_all"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    must_have_one_noun+=("all")
    must_have_one_noun+=("frames")
    must_have_one_noun+=("projects")
    must_have_one_noun+=("tags")
    noun_aliases=()
}

_tom_remove_frame()
{
    last_command="tom_remove_frame"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_remove_project()
{
    last_command="tom_remove_project"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--name-delimiter=")
    two_word_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_remove()
{
    last_command="tom_remove"

    command_aliases=()

    commands=()
    commands+=("all")
    commands+=("frame")
    commands+=("project")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_rename()
{
    last_command="tom_rename"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_report()
{
    last_command="tom_report"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--css-file=")
    two_word_flags+=("--css-file")
    flags_with_completion+=("--css-file")
    flags_completion+=("__tom_handle_filename_extension_flag gohtml")
    local_nonpersistent_flags+=("--css-file")
    local_nonpersistent_flags+=("--css-file=")
    flags+=("--day=")
    two_word_flags+=("--day")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--day")
    local_nonpersistent_flags+=("--day=")
    local_nonpersistent_flags+=("-d")
    flags+=("--decimal")
    local_nonpersistent_flags+=("--decimal")
    flags+=("--description=")
    two_word_flags+=("--description")
    local_nonpersistent_flags+=("--description")
    local_nonpersistent_flags+=("--description=")
    flags+=("--from=")
    two_word_flags+=("--from")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--from")
    local_nonpersistent_flags+=("--from=")
    local_nonpersistent_flags+=("-f")
    flags+=("--include-archived")
    local_nonpersistent_flags+=("--include-archived")
    flags+=("--json")
    local_nonpersistent_flags+=("--json")
    flags+=("--matrix-tables")
    local_nonpersistent_flags+=("--matrix-tables")
    flags+=("--month")
    flags+=("-m")
    local_nonpersistent_flags+=("--month")
    local_nonpersistent_flags+=("-m")
    flags+=("--output-file=")
    two_word_flags+=("--output-file")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output-file")
    local_nonpersistent_flags+=("--output-file=")
    local_nonpersistent_flags+=("-o")
    flags+=("--project=")
    two_word_flags+=("--project")
    two_word_flags+=("-p")
    local_nonpersistent_flags+=("--project")
    local_nonpersistent_flags+=("--project=")
    local_nonpersistent_flags+=("-p")
    flags+=("--project-delimiter=")
    two_word_flags+=("--project-delimiter")
    local_nonpersistent_flags+=("--project-delimiter")
    local_nonpersistent_flags+=("--project-delimiter=")
    flags+=("--round-frames=")
    two_word_flags+=("--round-frames")
    local_nonpersistent_flags+=("--round-frames")
    local_nonpersistent_flags+=("--round-frames=")
    flags+=("--round-frames-to=")
    two_word_flags+=("--round-frames-to")
    local_nonpersistent_flags+=("--round-frames-to")
    local_nonpersistent_flags+=("--round-frames-to=")
    flags+=("--save-config=")
    two_word_flags+=("--save-config")
    local_nonpersistent_flags+=("--save-config")
    local_nonpersistent_flags+=("--save-config=")
    flags+=("--short-titles")
    local_nonpersistent_flags+=("--short-titles")
    flags+=("--show-empty")
    local_nonpersistent_flags+=("--show-empty")
    flags+=("--show-sales")
    local_nonpersistent_flags+=("--show-sales")
    flags+=("--show-stop-time")
    local_nonpersistent_flags+=("--show-stop-time")
    flags+=("--show-summary")
    local_nonpersistent_flags+=("--show-summary")
    flags+=("--show-tracked")
    local_nonpersistent_flags+=("--show-tracked")
    flags+=("--show-untracked")
    local_nonpersistent_flags+=("--show-untracked")
    flags+=("--split=")
    two_word_flags+=("--split")
    two_word_flags+=("-s")
    local_nonpersistent_flags+=("--split")
    local_nonpersistent_flags+=("--split=")
    local_nonpersistent_flags+=("-s")
    flags+=("--subprojects")
    local_nonpersistent_flags+=("--subprojects")
    flags+=("--template=")
    two_word_flags+=("--template")
    local_nonpersistent_flags+=("--template")
    local_nonpersistent_flags+=("--template=")
    flags+=("--template-file=")
    two_word_flags+=("--template-file")
    flags_with_completion+=("--template-file")
    flags_completion+=("__tom_handle_filename_extension_flag gohtml")
    local_nonpersistent_flags+=("--template-file")
    local_nonpersistent_flags+=("--template-file=")
    flags+=("--title=")
    two_word_flags+=("--title")
    local_nonpersistent_flags+=("--title")
    local_nonpersistent_flags+=("--title=")
    flags+=("--to=")
    two_word_flags+=("--to")
    two_word_flags+=("-t")
    local_nonpersistent_flags+=("--to")
    local_nonpersistent_flags+=("--to=")
    local_nonpersistent_flags+=("-t")
    flags+=("--week=")
    two_word_flags+=("--week")
    two_word_flags+=("-w")
    local_nonpersistent_flags+=("--week")
    local_nonpersistent_flags+=("--week=")
    local_nonpersistent_flags+=("-w")
    flags+=("--year=")
    two_word_flags+=("--year")
    two_word_flags+=("-y")
    local_nonpersistent_flags+=("--year")
    local_nonpersistent_flags+=("--year=")
    local_nonpersistent_flags+=("-y")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_start()
{
    last_command="tom_start"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--create-missing")
    local_nonpersistent_flags+=("--create-missing")
    flags+=("--notes=")
    two_word_flags+=("--notes")
    local_nonpersistent_flags+=("--notes")
    local_nonpersistent_flags+=("--notes=")
    flags+=("--stop-on-start")
    local_nonpersistent_flags+=("--stop-on-start")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_status_projects()
{
    last_command="tom_status_projects"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--archived")
    local_nonpersistent_flags+=("--archived")
    flags+=("--delimiter=")
    two_word_flags+=("--delimiter")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--delimiter")
    local_nonpersistent_flags+=("--delimiter=")
    local_nonpersistent_flags+=("-d")
    flags+=("--format=")
    two_word_flags+=("--format")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    local_nonpersistent_flags+=("-f")
    flags+=("--include-active")
    local_nonpersistent_flags+=("--include-active")
    flags+=("--name-delimiter=")
    two_word_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--output=")
    two_word_flags+=("--output")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output")
    local_nonpersistent_flags+=("--output=")
    local_nonpersistent_flags+=("-o")
    flags+=("--show-empty")
    flags+=("-e")
    local_nonpersistent_flags+=("--show-empty")
    local_nonpersistent_flags+=("-e")
    flags+=("--show-overall")
    local_nonpersistent_flags+=("--show-overall")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_status()
{
    last_command="tom_status"

    command_aliases=()

    commands=()
    commands+=("projects")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--delimiter=")
    two_word_flags+=("--delimiter")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--delimiter")
    local_nonpersistent_flags+=("--delimiter=")
    local_nonpersistent_flags+=("-d")
    flags+=("--format=")
    two_word_flags+=("--format")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    local_nonpersistent_flags+=("-f")
    flags+=("--name-delimiter=")
    two_word_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter")
    local_nonpersistent_flags+=("--name-delimiter=")
    flags+=("--verbose")
    flags+=("-v")
    local_nonpersistent_flags+=("--verbose")
    local_nonpersistent_flags+=("-v")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_stop()
{
    last_command="tom_stop"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--all")
    flags+=("-a")
    local_nonpersistent_flags+=("--all")
    local_nonpersistent_flags+=("-a")
    flags+=("--notes=")
    two_word_flags+=("--notes")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--notes")
    local_nonpersistent_flags+=("--notes=")
    local_nonpersistent_flags+=("-n")
    flags+=("--past=")
    two_word_flags+=("--past")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--past")
    local_nonpersistent_flags+=("--past=")
    local_nonpersistent_flags+=("-d")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_tags()
{
    last_command="tom_tags"

    command_aliases=()

    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--delimiter=")
    two_word_flags+=("--delimiter")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--delimiter")
    local_nonpersistent_flags+=("--delimiter=")
    local_nonpersistent_flags+=("-d")
    flags+=("--format=")
    two_word_flags+=("--format")
    two_word_flags+=("-f")
    local_nonpersistent_flags+=("--format")
    local_nonpersistent_flags+=("--format=")
    local_nonpersistent_flags+=("-f")
    flags+=("--output=")
    two_word_flags+=("--output")
    two_word_flags+=("-o")
    local_nonpersistent_flags+=("--output")
    local_nonpersistent_flags+=("--output=")
    local_nonpersistent_flags+=("-o")
    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_tom_root_command()
{
    last_command="tom"

    command_aliases=()

    commands=()
    commands+=("cancel")
    commands+=("config")
    commands+=("create")
    commands+=("edit")
    commands+=("frames")
    commands+=("help")
    commands+=("import")
    commands+=("projects")
    commands+=("remove")
    commands+=("rename")
    commands+=("report")
    commands+=("start")
    commands+=("status")
    commands+=("stop")
    commands+=("tags")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--backup-dir=")
    two_word_flags+=("--backup-dir")
    flags+=("--config=")
    two_word_flags+=("--config")
    two_word_flags+=("-c")
    flags+=("--data-dir=")
    two_word_flags+=("--data-dir")
    flags+=("--iso-dates")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_tom()
{
    local cur prev words cword split
    declare -A flaghash 2>/dev/null || :
    declare -A aliashash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __tom_init_completion -n "=" || return
    fi

    local c=0
    local flag_parsing_disabled=
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("tom")
    local command_aliases=()
    local must_have_one_flag=()
    local must_have_one_noun=()
    local has_completion_function=""
    local last_command=""
    local nouns=()
    local noun_aliases=()

    __tom_handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_tom tom
else
    complete -o default -o nospace -F __start_tom tom
fi

# ex: ts=4 sw=4 et filetype=sh
