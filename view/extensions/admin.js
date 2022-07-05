
export * from './tags/js/admin';
import * as flarum_mentions from './mentions/js/admin';
export * from './likes/js/admin';
export * from './custom-footer/js/admin';
// export * from './auth-github/js/admin';
// export * from './analytics/js/admin';
flarum.extensions["flarum-mentions"] = flarum_mentions;
