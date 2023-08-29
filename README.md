<h1>News Feed Bot</h1>

<h2>Features</h2>
    <ul>
        <li>Fetching articles from RSS feeds</li>
        <li>Article summaries powered by GPT-3.5</li>
        <li>Admin commands for managing sources</li>
    </ul>

<h2>Configuration</h2>

<h3>Environment variables</h3>
    <ul>
        <li><code>NFB_TELEGRAM_BOT_TOKEN</code> — token for Telegram Bot API</li>
        <li><code>NFB_TELEGRAM_CHANNEL_ID</code> — ID of the channel to post to, can be obtained via @JsonDumpBot</li>
        <li><code>NFB_DATABASE_DSN</code> — PostgreSQL connection string</li>
        <li><code>NFB_FETCH_INTERVAL</code> — the interval of checking for new articles, default 10m</li>
        <li><code>NFB_NOTIFICATION_INTERVAL</code> — the interval of delivering new articles to Telegram channel, default 1m</li>
        <li><code>NFB_OPENAI_KEY</code> — token for OpenAI API</li>
        <li><code>NFB_OPENAI_PROMPT</code> — prompt for GPT-3.5 Turbo to generate a summary</li>
    </ul>

<h3>HCL</h3>
    <p>News Feed Bot can be configured with HCL config file. The service is looking for a config file in the following locations:</p>
    <ul>
        <li>./config.hcl</li>
        <li>./config.local.hcl</li>
        <li>$HOME/.config/news-feed-bot/config.hcl</li>
    </ul>
    <p>The names of parameters are the same except that there is no prefix, and names are in lowercase instead of uppercase.</p>

<h2>Nice to have features (backlog)</h2>
    <ul>
        <li>More types of resources — not only RSS</li>
        <li>Summary for the article</li>
    </ul>
