#!/usr/bin/env bash
# set -x
dst_dir="results"
template="template.cpp"
target_symbol="__HOGEHOGE__"
clang_lib_path='-I/usr/local/Cellar/llvm/6.0.0/include'
function xargs-printf() {
	while read line || [ -n "${line}" ]; do
		printf "$@" $line
	done
}

[[ ! -e $dst_dir ]] && mkdir -p $dst_dir
cd $dst_dir

(cd /usr/local/Cellar/llvm/6.0.0/include/clang && find . -name "*.h" | xargs-printf '#include "clang/%s"\n') >$template
cat >>$template <<EOF
void dummy(){
$target_symbol
}
EOF

declare -A map
declare -A map_tmp
function parse() {
	target=$1
	comp_target="$1::"
	namespace=${target##*::}
	echo $namespace
	filename=$(echo "$target" | sed 's/:/-/g').cpp
	cat $template | sed 's/'$target_symbol'/'$comp_target'/g' >"$filename"
	line_no=$(sed -n '/'$target_symbol'/=' $template)
	for target in $(clang++ -Xclang -fsyntax-only -Xclang -code-completion-at=$filename:$line_no:$((${#comp_target} + 1)) "$filename" -c $clang_lib_path | tee $filename.comp |
		grep -v '#' | awk '{if ($2 != "'$namespace'") print $2; }' | xargs-printf "$comp_target%s\n"); do
		# 		echo $target
		[[ ${map_tmp[$target]} == '' ]] && map_tmp[$target]='not yet'
	done
}

recursive_flag=0
base_target="clang"
[[ $# -eq 1 ]] && base_target="$1"
# base_target="clang::vfs::detail"
map[$base_target]='not yet'

echo "base_target:$base_target"

flag=1
while [[ $flag == 1 ]]; do
	for target in ${!map[*]}; do
		if [[ ${map[$target]} == 'not yet' ]]; then
			echo "parse $target"
			parse $target
			map[$target]='done'
		fi
	done

	[[ $recursive_flag == 0 ]] && break
	flag=0
	for target in ${!map_tmp[*]}; do
		[[ $(echo $target | fgrep $base_target) ]] || echo "skip $target" && continue
		[[ ${map[$target]} == '' ]] && map[$target]=${map_tmp[$target]} && echo "add parse list $target" && flag=1
	done
done
