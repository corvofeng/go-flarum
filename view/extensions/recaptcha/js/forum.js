import { extend, override } from 'flarum/extend';
import app from 'flarum/app';
import Component from 'flarum/Component';
import LogInButtons from 'flarum/components/LogInButtons';
import LogInButton from 'flarum/components/LogInButton';
import SignUpModal from 'flarum/components/SignUpModal';
import ModalManager from 'flarum/components/ModalManager';
import LoadingIndicator from 'flarum/components/LoadingIndicator';

export * from './src/forum';


// app.initializers.add('flarum-auth-github', () => {
//   extend(LogInButtons.prototype, 'items', function(items) {
//     items.add('github',
//       <LogInButton
//         className="Button LogInButton--github"
//         icon="fab fa-github"
//         path="/auth/github">
//         {app.translator.trans('flarum-auth-github.forum.log_in.with_github_button')}
//       </LogInButton>
//     );
//   });
// });
