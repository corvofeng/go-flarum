
export * from '../framework/extensions/emoji/js/forum';
import * as flarum_tags from '../framework/extensions/tags/js/forum';
export * from '../framework/extensions/markdown/js/forum';
import * as flarum_likes from '../framework/extensions/likes/js/forum';
import * as flarum_mentions from './mentions/js/forum';

export * from './recaptcha/js/forum';
import * as flarum_subscriptions from './flarum-subscriptions/js/forum';
// export * from './diff/js/forum';

export * from './auth-github/js/forum';
export * from './analytics/js/forum';
import * as flarum_sticky from './flarum-sticky/js/forum';
// export * from './custom-footer/js/forum';
// export * from './flarum-pipetables/js/forum';
flarum.extensions["flarum-sticky"] = flarum_sticky;
flarum.extensions["flarum-subscriptions"] = flarum_subscriptions;
flarum.extensions["flarum-mentions"] = flarum_mentions;
flarum.extensions["tags"] = flarum_tags;
flarum.extensions["likes"] = flarum_likes;