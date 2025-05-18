#!/usr/bin/env bash

# Files like this:
# ```
# + generator
#  - templates
#  - tto.conf
#  - example.sh
# ```

TABLE_PATTERN="*"
CURRENT_DIR=$(dirname "$0")
COMMON_DIR="${CURRENT_DIR}"/../maple-transport/src/main/java/com/ruoyi/maple
CONTROLLER_DIR="${CURRENT_DIR}"/../ruoyi-admin/src/main/java/com/ruoyi/web
VIEW_DIR="${CURRENT_DIR}"/../board/src/packages/transport/views
GO_DIR="${CURRENT_DIR}"/../tgbot/src

echo "表名关键字: $TABLE_PATTERN"
echo "核心目录: $COMMON_DIR"
echo "控制器目录: $CONTROLLER_DIR"
echo "页面目录: $VIEW_DIR"

echo "开始生成代码..."

if [[ $(whereis tto) = 'tto:' ]]; then
  echo '未安装tto客户端,请运行安装命令： curl -L https://raw.githubusercontent.com/ixre/tto/master/install | sh'
fi

CONF_DIR=$(dirname "$0")
tto  -conf "$CONF_DIR"/tto.conf -t "$CONF_DIR"/templates -o output \
-model '' \
-pkg '' \
-excludes tmp_ \
-table "${TABLE_PATTERN}" \
-clean


copy_dir() {
  find "${CURRENT_DIR}/output/${1}" -type d -name "$2" -exec cp -r {} "$3" \;
}
copy_java_modules(){
  copy_dir "spring" "entity" ${COMMON_DIR}
  copy_dir "spring" "mapper" ${COMMON_DIR}
  copy_dir "spring" "service" ${COMMON_DIR}
  copy_dir "spring" "controller" ${CONTROLLER_DIR}
}


echo '' && echo "请选择操作："
echo "----------------------------------------"
echo "【1】生成Java代码到项目中"
echo "【2】生成Java和Vue代码到项目中"
echo "【3】生成Go代码到项目中"
echo "----------------------------------------"

while true; do
  read -p "请输入操作编号(回车退出): " operation
  case $operation in
    '')break;;
    1)
      copy_java_modules
      break
      ;;
    2)
      copy_java_modules
      VIEW_OUTPUT_DIR=$(find "${CURRENT_DIR}"/output -type d -name "views")
      find "${VIEW_OUTPUT_DIR}" -maxdepth 1 -type d -not -name "views" -exec cp -r {} "${VIEW_DIR}" \;
      break
      ;;
    3)
      find "${CURRENT_DIR}/output/go" -maxdepth 1 -type d -not -name "go" -exec cp -r {} "${GO_DIR}" \;
      break
      ;;
    *)
      echo "无效的操作编号，请重新输入。"
      ;;
  esac
done


# Replace generator description part of code file
# find output/spring -name "*.java" -print0 |  xargs -0 sed -i ':label;N;s/This.*Copy/Copy/g;b label'
# Replace package
# find output/spring -name "*.java" -print0 |  xargs -0 sed -i 's/net.fze/com.pkg/g'
# Replace type
# find output/java -name "*.java"  -print0 | xargs -0 sed -i 's/ int / Integer /g'
# copy files to project folder
# find ./src -path "*/entity" -print0 | xargs -0 cp output/spring/src/main/java/com/github/tto/entity/*
