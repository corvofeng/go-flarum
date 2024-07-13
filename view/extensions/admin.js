

import * as flarum_tags from '../framework/extensions/tags/js/admin';
export * from '../framework/extensions/markdown/js/admin';
export * from '../framework/extensions/likes/js/admin';
import * as flarum_mentions from './mentions/js/admin';
import * as flarum_sticky from './flarum-sticky/js/admin';

// export * from './custom-footer/js/admin';
// export * from './auth-github/js/admin';
// export * from './analytics/js/admin';
flarum.extensions["flarum-sticky"] = flarum_sticky;
flarum.extensions["flarum-mentions"] = flarum_mentions;
flarum.extensions["tags"] = flarum_tags;
