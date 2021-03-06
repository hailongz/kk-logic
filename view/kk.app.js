
kk.Logic = function (app, name, object) {
    this.app = app;
    this.children = [];
    this.title = name + ': ' + (object['title'] || object['$class']);

    for (var key in object) {

        if (key.startsWith("on")) {

            var v = object[key];

            if (typeof v == 'string') {
                v = app.object[v];
            }

            if (typeof v == 'object' && v['$class']) {
                var logic = new kk.Logic(app, key.substr(2), v);
                this.children.push(logic);
            }
        }
    }
};

kk.Logic.prototype = kk.extend(Object.prototype, {

    appendTo: function (div) {

        var view = this.view = document.createElement("div");
        view.className = "kk-element";

        var span = document.createElement("span");
        span.innerText = this.title;

        view.appendChild(span);

        div.appendChild(view);

        this.each(function (logic) {
            logic.appendTo(div);
        });

    },

    viewSize: function () {
        return { width: this.view.clientWidth + 10, height: this.view.clientHeight };
    },

    each: function (fn) {
        for (var i = 0; i < this.children.length; i++) {
            var v = this.children[i];
            fn(v, i);
        }
    },

    moveTo: function (x, y, width) {

        var spacing = 10;

        this.view.style.left = x + 'px';
        this.view.style.top = y + 'px';

        if (width != 0) {
            this.view.style.width = width + 'px';
        }

        var maxWidth = this.getMaxWidth();
        var viewSize = this.viewSize();

        x = x + viewSize.width + spacing;

        this.each(function (logic) {
            logic.moveTo(x, y, maxWidth);
            var size = logic.size();
            y += size.height + spacing;
        });

    },

    getMaxWidth: function () {
        if (this.maxWidth === undefined) {
            var maxWidth = 0;
            this.each(function (logic) {
                var viewSize = logic.viewSize();
                if (viewSize.width > maxWidth) {
                    maxWidth = viewSize.width;
                }
            });
            this.maxWidth = maxWidth;
        }
        return this.maxWidth;
    },

    size: function () {

        var spacing = 10;

        if (this.width === undefined && this.height === undefined && this.view) {

            var maxWidth = 0;
            var dy = 0;

            this.each(function (logic) {
                var size = logic.size();
                dy += size.height + spacing;
                if (size.width > maxWidth) {
                    maxWidth = size.width;
                }
            });

            var viewSize = this.viewSize();

            if (maxWidth == 0) {
                this.width = viewSize.width;
                this.height = viewSize.height;
            } else {
                this.width = viewSize.width + maxWidth;
                this.height = dy ;
            }

        }

        return { width: this.width, height: this.height };
    }
});

kk.App = function (object) {
    this.object = object;

    var v = object['in'];

    if (typeof v == 'object') {
        this.in = new kk.Logic(this, 'in', v);
    }

};

kk.App.prototype = kk.extend(Object.prototype, {

    appendTo: function (div) {

        if (this.in) {

            var x = 40;
            var y = 40;

            this.in.appendTo(div);
            this.in.moveTo(x, y, 0);
            var size = this.in.size();
            div.style.width = (x + size.width) + 'px';
            div.style.height = (y + size.height) + 'px';

        }

    }
});
