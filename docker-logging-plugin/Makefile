PLUGIN_NAME = morawskim/docker-loggin-plugin
PLUGIN_TAG ?= dev
BUILD_ROOT ?= ~/plugin-build-root/

all: clean rootfs create

clean:
	@echo "### rm ${BUILD_ROOT}"
	@rm -rf ${BUILD_ROOT}

rootfs:
	@echo "### docker build"
	@docker build -q -t ${PLUGIN_NAME}:rootfs .
	@echo "### create rootfs directory in ./plugin/rootfs"
	@mkdir -p ${BUILD_ROOT}/rootfs
	@docker create --name tmp ${PLUGIN_NAME}:rootfs
	@docker export tmp | tar -x -C ${BUILD_ROOT}/rootfs
	@echo "### copy plugin configuration"
	cp plugin/config.json ${BUILD_ROOT}/config.json
	@docker rm -vf tmp

create:
	@echo "### remove existing plugin ${PLUGIN_NAME}:${PLUGIN_TAG} if exists"
	@docker plugin rm -f ${PLUGIN_NAME}:${PLUGIN_TAG} || true
	@echo "### create new plugin ${PLUGIN_NAME}:${PLUGIN_TAG} from ./plugin"
	@docker plugin create ${PLUGIN_NAME}:${PLUGIN_TAG} ${BUILD_ROOT}

enable:
	@echo "### enable plugin ${PLUGIN_NAME}:${PLUGIN_TAG}"
	@docker plugin enable ${PLUGIN_NAME}:${PLUGIN_TAG}
