<notification>
    <div>
        <div class="app">{opts.app}</div>
        <div class="timestamp">{opts.timestamp}</div>
        <div class="title">{opts.title}</div>
        <pre class="text">{opts.text}</pre>
    </div>

    <style>
        .app, .timestamp {
            font-weight: bold;
            color: #78B6B6;
        }
        .title {
            font-weight: bold;
            color: #3C8B8B;
        }
        .text {
            border-left: solid thin darkgray;
            padding-left: 2em;
        }
    </style>
</notification>
