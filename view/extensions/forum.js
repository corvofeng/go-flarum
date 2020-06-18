
export * from './emoji/js/forum';
export * from './recaptcha/js/forum';
import * as flarum_mentions from './mentions/js/forum';
export * from './tags/js/forum';
export * from './markdown/js/forum';

flarum.extensions["flarum-mentions"] = flarum_mentions;