$(function() {
    $.spinner_ext = {
        tmpl: '<div class="spinner-border" role="status">' +
            '<span class="visually-hidden">${Text}</span>' +
            '</div>',

        init: function(msg) {
            this.create(msg);
            this.set_css();

            return this.spinnerDiv;
        },

        template: function (tmpl, data) {
            $.each(data, function (k, v) {
                tmpl = tmpl.replace('${' + k + '}', v);
            });
            return $(tmpl);
        },

        create: function(msg) {
            this.spinnerDiv = this.template(this.tmpl, {
                Text: msg,
            }).hide();
            this.spinnerDiv.appendTo($('body'));
        },

        set_css: function () {
            var spinnerDiv = this.spinnerDiv;

            spinnerDiv.css({
                'position': 'fixed',
                'z-index': 10001 + $(".alert").length,
                'position': 'absolute',
                'top': '0',
                'left': '0'
            });
        },
    };

    $.spinner = function(msg) {
        return $.spinner_ext.init(msg);
    }
});
