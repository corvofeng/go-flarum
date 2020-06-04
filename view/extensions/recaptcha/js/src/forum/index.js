// import app from 'flarum/forum';


// console.log("Hello world");
// console.log("Hello");
// console.log(app);

import { extend } from 'flarum/extend';
import app from 'flarum/app';
import LogInButtons from 'flarum/components/LogInButtons';
import LogInButton from 'flarum/components/LogInButton';


  console.log("add github")
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
