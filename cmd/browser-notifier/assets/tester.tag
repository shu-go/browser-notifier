<tester>
    <span class="vline"></span>
    App: <input type="text" name="app" value={app} onkeyup={appChanged} />
    <span class="vline"></span>
    Title: <input type="text" name="title" value={title} onkeyup={titleChanged} />
    <span class="vline"></span>
    Text: <textarea name="text" rows="1" cols="30" onkeyup={textChanged}>{text}</textarea>
    <span class="vline"></span>
    <input type="button" value="Notify" onclick={notifyClicked} />
    =&gt; <textarea name="preview" rows="1" cols="30" readonly="readonly">{preview}</textarea>
    
    <style>
        .vline {
            border-left: inset thin black;
            width: 2px;
            margin-left: 1ex;
            margin-right: 1ex;
        }
    </style>

    // script

    this.app = opts.app;
    this.title = "";
    this.text = "";
    this.preview = "";

    appChanged(e) {
        this.app = e.target.value;
        this.updatePreview();
    }

    titleChanged(e) {
        this.title = e.target.value;
        this.updatePreview();
    }

    textChanged(e) {
        this.text = e.target.value;
        this.updatePreview();
    }

    updatePreview() {
        this.preview = this.buildJSONString();
    }

    buildJSONString() {
        var n = {"app": this.app, "title": this.title, "text": this.text};
        return JSON.stringify(n);
    }

    notifyClicked(e) {
        var that = this;

        var request = new XMLHttpRequest();
        request.open("post", "/notifications", true);
        request.onload = function (event) {
            that.title = "";
            that.text = "";
            that.updatePreview();
        };
        request.onerror = function (event) {
            console.log(event.type); // => "error"
        };
        request.send(this.buildJSONString());            
    }

</tester>
