import { extend } from 'flarum/common/extend';
import app from 'flarum/forum/app';
import LogInButton from 'flarum/forum/components/LogInButton';
import LogInButtons from 'flarum/forum/components/LogInButtons';

app.initializers.add('flarum-auth-github', () => {
  extend(LogInButtons.prototype, 'items', function(items) {
    items.add('github',
      <LogInButton
        className="Button LogInButton--github"
        icon="fab fa-github"
        path="/auth/github">
        {app.translator.trans('flarum-auth-github.forum.log_in.with_github_button')}
      </LogInButton>
    );
  });
});
