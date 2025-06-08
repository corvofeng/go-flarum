import app from 'flarum/app';

import GithubSettingsModal from './components/GithubSettingsModal';

app.initializers.add('flarum-auth-github', () => {
  app.extensionSettings['flarum-auth-github'] = () => app.modal.show(new GithubSettingsModal());
});
