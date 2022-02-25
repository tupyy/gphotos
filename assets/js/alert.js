(function($) {
    $.alert_ext = {
        defaults: {
            autoClose: true,  
            closeTime: 5000,   
            withTime: false, 
            type: 'danger',  
            position: ['center', [-0.42, 0]], 
            title: false, 
            close: '',   
            speed: 'normal',  
            isOnly: true, 
            minTop: 10, 
            onShow: function () {
            },  
            onClose: function () {
            }  
        },

        tmpl: '<div class="alert alert-dismissable fade show ${State}">' +
                    '<p>${Content}</p>' +
                    '<button type="button" class="btn-close close" data-bs-dismiss="alert" aria-hidden="true"></button>' +
                '</div>',

        init: function (msg, options) {
            this.options = $.extend({}, this.defaults, options);

            this.create(msg);
            this.set_css();

            this.bind_event();

            return this.alertDiv;
        },

        template: function (tmpl, data) {
            $.each(data, function (k, v) {
                tmpl = tmpl.replace('${' + k + '}', v);
            });
            return $(tmpl);
        },

        create: function (msg) {
            this.alertDiv = this.template(this.tmpl, {
                State: 'alert-' + this.options.type,
                Content: msg
            }).hide();
            if (this.options.isOnly) {
                $('body > .alert').remove();
            }
            this.alertDiv.appendTo($('body'));
        },

        set_css: function () {
            var alertDiv = this.alertDiv;

            alertDiv.css({
                'position': 'fixed',
                'z-index': 10001 + $(".alert").length,
                'display': 'flex'
            });
            
            $('p', this.alertDiv).css({
                'margin-right': '15px',
                'margin-bottom': '0',
            });

            var ie6 = 0;
            if ($.browser && $.browser.msie && $.browser.version == '6.0') {
                alertDiv.css('position', 'absolute');
                ie6 = $(window).scrollTop();
            }

            var position = this.options.position,
                pos_str = position[0].split('-'),
                pos = [0, 0];
            if (position.length > 1) {
                pos = position[1];
            }

            if (pos[0] > -1 && pos[0] < 1) {
                pos[0] = pos[0] * $(window).height();
            }
            if (pos[1] > -1 && pos[1] < 1) {
                pos[1] = pos[1] * $(window).width();
            }


            for (var i in pos_str) {
                if ($.type(pos_str[i]) !== 'string') {
                    continue;
                }
                var str = pos_str[i].toLowerCase();

                if ($.inArray(str, ['left', 'right']) > -1) {
                    alertDiv.css(str, pos[1]);
                } else if ($.inArray(str, ['top', 'bottom']) > -1) {
                    alertDiv.css(str, pos[0] + ie6);
                } else {
                    alertDiv.css({
                        'top': ($(window).height() - alertDiv.outerHeight()) / 2 + pos[0] + ie6,
                        'left': ($(window).width() - alertDiv.outerWidth()) / 2 + pos[1]
                    });
                }
            }

            if (parseInt(alertDiv.css('top')) < this.options.minTop) {
                alertDiv.css('top', this.options.minTop);
            }
        },

        bind_event: function () {
            this.bind_show();
            this.bind_close();

            if ($.browser && $.browser.msie && $.browser.version == '6.0') {
                this.bind_scroll();
            }
        },

        bind_show: function () {
            var ops = this.options;
            this.alertDiv.fadeIn(ops.speed, function () {
                ops.onShow($(this));
            });
        },

        bind_close: function () {
            var alertDiv = this.alertDiv,
                ops = this.options,
                closeBtn = $('.close', alertDiv).add($(this.options.close, alertDiv));

            closeBtn.bind('click', function (e) {
                alertDiv.fadeOut(ops.speed, function () {
                    $(this).remove();
                    ops.onClose($(this));
                });
                e.stopPropagation();
            });

            if (this.options.autoClose) {
                var time = parseInt(this.options.closeTime / 1000);
                if (this.options.withTime) {
                    $('p', alertDiv).append('<span>...<em>' + time + '</em></span>');
                }
                var timer = setInterval(function () {
                    $('em', alertDiv).text(--time);
                    if (!time) {
                        clearInterval(timer);
                        closeBtn.trigger('click');
                    }
                }, 1000);
            }
        },

        bind_scroll: function () {
            var alertDiv = this.alertDiv,
                top = alertDiv.offset().top - $(window).scrollTop();
            $(window).scroll(function () {
                alertDiv.css("top", top + $(window).scrollTop());
            })
        },

        check_mobile: function () {
            var userAgent = navigator.userAgent;
            var keywords = ['Android', 'iPhone', 'iPod', 'iPad', 'Windows Phone', 'MQQBrowser'];
            for (var i in keywords) {
                if (userAgent.indexOf(keywords[i]) > -1) {
                    return keywords[i];
                }
            }
            return false;
        }
    };

    $.alert = function (msg, arg) {
        if ($.alert_ext.check_mobile()) {
            alert(msg);
            return;
        }
        if (!$.trim(msg).length) {
            return false;
        }
        if ($.type(arg) === "string") {
            arg = {
                title: arg
            }
        }
        if (arg && arg.type == 'error') {
            arg.type = 'danger';
        }
        return $.alert_ext.init(msg, arg);
    }
})(jQuery);
