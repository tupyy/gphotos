$(function () {
    $.widget("upload.fileui", $.upload.fileUISkin, {
        options: {
            signed_url: undefined,
            file: undefined,
            filename: undefined,
            id: undefined,
            progress: 0,
            uploaded: false,
            xhr: undefined,
            onFinishUploadCallback: undefined
        },
        _create: function () {
            this._super();
            this.options.id = this.guid();
        },
        _delete: function () {
            this._trigger("delete", event, {id: this.options.id});
        },

        _setOption: function (key, value) {
            this._super(key, value);
            this.refresh();
        },

        _setOptions: function (options) {
            this._super(options);
            this.refresh();
        },

        setUrl: function (url) {
            this._initXHRData(url);
        },
        setCallback: function(callback) {
            this.options.onFinishedUpload = callback;
        },
        _initXHRData: function (url) {
            const postData = new FormData();
            postData.append('file', this.options.file);
            this.options.data = postData;
            this.options.processData = false;
            this.options.type = 'POST';
            this.options.url = url;
        },
        send: function () {
            this.options.xhr = new XMLHttpRequest();
            let xhr = this.options.xhr;
            let dfd = $.Deferred();
            
            // call callback when resolved
            dfd.done(() => {
                this.options.onFinishedUpload(this.options.id);
            });

            if (this.options.uploaded) {
                dfd.resolve();
            }

            // set status
            this.updateStatus("Uploading");

            const self = this;
            xhr.upload.addEventListener("progress", function (e) {
                if (e.lengthComputable) {
                    const progress = Math.round((e.loaded * 100) / e.total);
                    self.refresh(progress);
                }
            }, false);

            xhr.onreadystatechange = function (e) {
                if (xhr.readyState === 4) {
                    self.refresh(100);
                    self.options.uploaded = true;
                    self.options.xhr = undefined;
                    if (xhr.status === 200 || xhr.status === 204) {
                        self.uploadFinished();
                        dfd.resolve();
                    } else {
                        dfd.fail(self.options.id);
                        self.updateStatus("Failed");
                    }
                }
            };

            xhr.open(this.options.type, this.options.url, true);
            xhr.send(this.options.data);

            let promise = dfd.promise();
            promise.abort = this.abort;
            return promise;
        },

        abort: function () {
            if (this.options.xhr) {
                this.options.xhr.abort();
                this.options.xhr = undefined;
            }
        },

        isUploading: function () {
            return this.options.xhr !== undefined;
        },

        guid: function () {
            function _p8(s) {
                var p = (Math.random().toString(16) + "000000000").substr(2, 8);
                return s ? "-" + p.substr(0, 4) + "-" + p.substr(4, 4) : p;
            }

            return _p8() + _p8(true) + _p8(true) + _p8();
        }
        ,
        destroy: function () {
            this.abort();
            this._off(this.element.find('.button'), 'click');
            this.element.remove();
        }
    })
    ;
});
