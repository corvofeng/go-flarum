import { extend, override } from 'flarum/extend';
import SignUpModal from 'flarum/components/SignUpModal';
import LogInModal from 'flarum/components/LogInModal';
import LogInButtons from 'flarum/components/LogInButtons';
import LogInButton from 'flarum/components/LogInButton';
import Alert from 'flarum/components/Alert';
import md5 from 'md5';
import SelfCaptcha from "./components/captcha";
import Captcha from "./captcha";

function makeAlert(type, content) {
  return {
    type: type,
    content: content,
  };
}

app.initializers.add('recaptcha', () => {
  let isOverrideSignUp = false;  // 是否重写了重要函数
  let isOverrideLogIn = false;  // 是否重写了重要函数

  // 处理注册页面
  extend(SignUpModal.prototype, 'submitData', function (data) {
    const recaptcha = this.fields().items.recaptcha;
    if (recaptcha) {
      data['captcha-solution'] = Captcha.captchaSolution;
      data['captcha-id'] = Captcha.captchaId;
    }
  });

  extend(SignUpModal.prototype, 'fields', function (fields) {
    if (!isOverrideSignUp) {
      isOverrideSignUp = true;
      override(SignUpModal.prototype, 'onsubmit', function (original, e) {
        e.preventDefault();
        this.loading = true;
        const data = this.submitData();
        const that = this;
        app.request({
          url: app.forum.attribute('baseUrl') + '/register',
          method: 'POST',
          body: data,
        }).then(
          (data) => {
            if (data.retcode !== 200) {
              Captcha.onrefresh();
              that.alertAttrs = makeAlert('error', data.retmsg);
            } else {
              that.alertAttrs = makeAlert('success', data.retmsg + "3秒后跳转到登录")
              setTimeout(() => {
                that.logIn();
              }, 3000);
            }

            that.loading = false;
            m.redraw();
          }
        );
      });
    }
    fields.add('recaptcha', <SelfCaptcha refreshBinder={Captcha.bindRefreshCallback}></SelfCaptcha>, 5);
  });

  let recaptcha = new SelfCaptcha();
  // 处理登录页面
  extend(LogInModal.prototype, 'fields', function (fields) {
    if (!isOverrideLogIn) {
      isOverrideLogIn = true;

      override(LogInModal.prototype, 'onsubmit', function (original, e) {
        e.preventDefault();
        this.loading = true;

        const identification = this.identification();
        const password = this.password();
        const remember = this.remember();
        const captchaId = Captcha.captchaId;
        const captchaSolution = Captcha.captchaSolution;

        const that = this;
        app.session
          .login({
            identification,
            'password': md5(password),
            remember,
            'captcha-solution': captchaSolution,
            'captcha-id': captchaId,
          }, { errorHandler: this.onerror.bind(this) })
          .then((data) => {
            if (data.retcode !== 200) {
              Captcha.onrefresh();
              that.alertAttrs = makeAlert('error', data.retmsg);
            } else {
              that.alertAttrs = makeAlert('success', data.retmsg + "3秒后刷新页面");
              setTimeout(() => {
                window.location.reload();
              }, 3000);
            }
            this.loading = false;
            m.redraw();
          });
      });
    }
    fields.add('recaptcha', <SelfCaptcha refreshBinder={Captcha.bindRefreshCallback}></SelfCaptcha>, 5);
  });

});