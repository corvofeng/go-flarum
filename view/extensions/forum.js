
export * from './emoji/js/forum';
export * from './recaptcha/js/forum';
import * as flarum_mentions from './mentions/js/forum';
export * from './tags/js/forum';
export * from './markdown/js/forum';
export * from './diff/js/forum';
export * from './likes/js/forum';
export * from './auth-github/js/forum';
export * from './analytics/js/forum';
export * from './flarum-pipetables/js/forum';

flarum.extensions["flarum-mentions"] = flarum_mentions;