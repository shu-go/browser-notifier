<notification-list>
    <div class="header fixed">
        <div class="fixbar">
            <span style="font-weight:bold;">App Filter</span> <input type="text" onkeyup={filterEdited} value={filter}></input>
            <ul class="apps" style="display: inline-block">
                <li class={selected : filter==""} onclick={filterChosen}>(all)</li>
                <li each={value in apps} class={selected : value.indexOf(filter)!=-1} onclick={filterChosen}>{value}</li>
            </ul>
            <div>
                <label for="tester_state">
                    <input id="tester_state" type="checkbox" value="checked" onclick={testerToggling} />Tester
                </label>
                <tester show={showTester} app="tester" />
            </div>
        </div>
    </div>

    <h1>Timeline</h1>
    <ul>
        <notification each={notifications.filter(appFilter)} app={app} title={title} text={text} timestamp={timestamp} />
    </ul>
    <div class="finishline"></div>

    <style>
        .header {
            height: 100px;
        }
        .header.fixed .fixbar {
            height: 100px;
            padding-left: 1em;

            background-color: #C5E3E3;
            position: fixed;
            top: 0px;

            width: 100%;
            margin: 0;
        }

        ul.apps li {
            display: inline-block;
            padding: 3px 5px;
            margin: 0px 3px;

            border: solid medium white;
            color: #3B6A83;
            font-weight: bold;
        }
        ul.apps li.selected {
            background-color: #1E6969;
            color: white;
        }

        .finishline {
            background: none;
            width: 100%;
            height: 20px;
        }
        .finishline.shown {
            background: linear-gradient(#F5FCFC, #C5E3E3);
            width: 100%;
            height: 20px;
        }
    </style>

    // script

    this.ws = null;

    this.notifications = opts.notifications;
    this.apps = [];
    this.filter = "";
    this.showTester = false;

    load(e) {
        var that = this;

        var request = new XMLHttpRequest();
        request.open("get", "/notifications", true);
        request.onload = function (event) {
            if (request.readyState === 4) {
                if (request.status === 200) {
                    that.notifications = JSON.parse(request.responseText);
                    that.accumerateApps();
                    that.update();

                } else {
                    console.log(request.statusText); // => Error Message
                }
            }
        };
        request.onerror = function (event) {
            console.log(event.type); // => "error"
        };
        request.send(null);            
    }

    accumerateApps() {
        var dict = {}
        var ns = this.notifications;
        for (idx in this.notifications) {
            if (ns[idx].app != "") {
                dict[ns[idx].app] = 1;
            }
        }
        var list = [];
        for (key in dict) {
            list.push(key);
        }
        this.apps = list.sort();
    }

    appFilter(n) {
        var f = this.filter;
        if (f == null || f == "" || n.app.indexOf(this.filter) != -1)
            return n
    }

    filterEdited(e) {
        this.filter = e.target.value;
    }

    filterChosen(e) {
        if (e.item) {
            this.filter = e.item.value;
        } else {
            this.filter = "";
        }
        return false;
    }

    testerToggling(e) {
        this.showTester = e.target.checked;
    }

    this.on('mount', function(){
        this.load(null);

        Notification.requestPermission();

        var that = this;
        var ws = new WebSocket("ws://"+window.location.host+"/push");
        this.ws = ws;
        ws.onerror = function(e) {
            console.log(e);
        };
        ws.onmessage = function(e) {
            var n = JSON.parse(e.data)
            nopts = {body: n.text, icon: "/favicon.png"}
            var bn = new Notification(n.title + (n.app != "" ? " (" + n.app + ")" : ""), nopts)
            bn.onclick = function(e){ this.close(); window.focus(); }

            that.load(null);
        };
    });

    this.on('updated', function() {
        window.scrollTo(0,document.body.scrollHeight);

        var fl = document.querySelector(".finishline");
        var body = document.body;
        if (body.scrollHeight > body.clientHeight) {
            fl.classList.add("shown");
        } else {
            fl.classList.remove("shown");
        }
    });

    window.addEventListener('unload', function(e) {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
    });
</notification-list>
