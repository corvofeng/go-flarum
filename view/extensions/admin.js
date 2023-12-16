

export * as flarum_tags from '../framework/extensions/tags/js/admin';
export * from '../framework/extensions/markdown/js/admin';
export * from '../framework/extensions/likes/js/admin';
import * as flarum_mentions from './mentions/js/admin';

// export * from './custom-footer/js/admin';
// export * from './auth-github/js/admin';
// export * from './analytics/js/admin';
flarum.extensions["flarum-mentions"] = flarum_mentions;
flarum.extensions["flarum-tags"] = flarum_tags;
