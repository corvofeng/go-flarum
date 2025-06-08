
import * as flarum_emoji from '../framework/extensions/emoji/js/forum';
flarum.extensions["flarum-emoji"] = flarum_emoji;

import * as flarum_tags from '../framework/extensions/tags/js/forum';
import * as flarum_markdown from '../framework/extensions/markdown/js/forum';
// import * as flarum_mentions from './mentions/js/forum';

import * as flarum_recaptcha from './recaptcha/js/forum';
import * as flarum_subscriptions from '../framework/extensions/subscriptions/js/forum';
import * as flarum_mentions from '../framework/extensions/mentions/js/forum';
// export * from './diff/js/forum';

// import * as flarum_fof_oauth from './fof-oauth/js/forum';
import * as flarum_github from './auth-github/js/forum';
export * from './analytics/js/forum';
// export * from './custom-footer/js/forum';
import * as flarum_likes from '../framework/extensions/likes/js/forum';
import * as flarum_sticky from '../framework/extensions/sticky/js/forum';

// export * from './flarum-pipetables/js/forum';
// flarum.extensions["flarum-fof-auth"] = flarum_fof_oauth;
flarum.extensions["flarum-github"] = flarum_github;
flarum.extensions["flarum-recaptch"] = flarum_recaptcha;
flarum.extensions["flarum-subscriptions"] = flarum_subscriptions;
flarum.extensions["flarum-markdown"] = flarum_markdown;
flarum.extensions["flarum-mentions"] = flarum_mentions;
flarum.extensions["tags"] = flarum_tags;
flarum.extensions["flarum-sticky"] = flarum_sticky;
flarum.extensions["likes"] = flarum_likes;