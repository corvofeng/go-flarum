
FLARUM_DIR = ./view/flarum/
LESSC = ./node_modules/.bin/lessc

#	@yarn add components-font-awesome less

webbundle:
	@cat ${FLARUM_DIR}/less/common/variables.less \
		${FLARUM_DIR}/less/common/mixins.less \
		${FLARUM_DIR}/less/forum.less > /tmp/forum.less
	@${LESSC} \
		--include-path=node_modules/components-font-awesome/less/:node_modules/bootstrap/less:${FLARUM_DIR}/less:${FLARUM_DIR}less/common /tmp/forum.less static/css/flaurm/forum.css

	@cat ${FLARUM_DIR}/less/common/variables.less \
		${FLARUM_DIR}/less/common/mixins.less \
		${FLARUM_DIR}/less/admin.less > /tmp/admin.less
	@${LESSC} \
		--include-path=node_modules/components-font-awesome/less/:node_modules/bootstrap/less:${FLARUM_DIR}/less:${FLARUM_DIR}less/common /tmp/admin.less static/css/flaurm/admin.css

