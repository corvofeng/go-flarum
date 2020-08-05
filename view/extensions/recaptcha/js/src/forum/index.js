import { extend, override } from 'flarum/extend';
import Component from 'flarum/Component';
import SignUpModal from 'flarum/components/SignUpModal';
import LogInModal from 'flarum/components/LogInModal';
import LogInButtons from 'flarum/components/LogInButtons';
import LogInButton from 'flarum/components/LogInButton';
import Alert from 'flarum/components/Alert';
import md5 from 'md5';

// app.initializers.add('saleksin-auth-google', () => {
//   extend(LogInButtons.prototype, 'items', function (items) {
//     items.add('google',
//       <LogInButton
//         className="Button LogInButton--google"
//         icon="fab fa-google"
//         path="/auth/google">
//         {app.translator.trans('saleksin-auth-google.forum.log_in.with_google_button')}
//       </LogInButton>
//     );
//   });
// });

// 如果想使用Google的验证码, 参考这个代码
// https://github.com/FriendsOfFlarum/recaptcha/blob/master/js/src/forum/components/Recaptcha.js
class SelfCaptcha extends Component {
  init() {
    super.init();
    this.recaptchaSolution = m.prop(this.props.recaptcha || '');
    this.isLoadCaptcha = false;
    this.captchaId = "";

    this.isInit = false;
  }

  isProvided(field) {
    return this.props.provided && this.props.provided.indexOf(field) !== -1;

  }

  view() {
    if (!this.isInit) {
      this.isInit = true;
      this.refreshCaptcha();
    }
    return (
      <div className="Form-group">
        <img
          src={this.getUrl(this.captchaId)}
          style={{ display: this.isLoadCaptcha ? '' : 'none' }}
          onclick={this.refreshCaptcha.bind(this)}
        />
        <div style={{ display: this.isLoadCaptcha ? 'none' : '' }} >
          <i class="fa fa-spinner fa-spin"></i>Loading
        </div>

        <input className="FormControl" name="recaptcha"
          type="text" placeholder="验证码" value={this.recaptchaSolution()}
          onchange={m.withAttr('value', this.recaptchaSolution)}
        />
      </div>
    );
  }

  getUrl(cid) {
    if (cid === "") return "";
    return `/captcha/${cid}.png`
  }
  setCaptchaID(cid) {
    this.captchaId = cid;
    this.isLoadCaptcha = true;
    m.redraw();
  }

  refreshCaptcha() {
    return app.request({
      url: app.forum.attribute('apiUrl') + '/new_captcha',
      method: 'GET',
    }).then((data) => {
      this.setCaptchaID(data['newCaptchaID'])
    });
  }

  getSolution() {
    return this.recaptchaSolution()
  }
  getCaptchaID() {
    return this.captchaId;
  }
}

app.initializers.add('recaptcha', () => {
  let recaptcha = new SelfCaptcha();
  let isOverrideSignUp = false;  // 是否重写了重要函数
  let isOverrideLogIn = false;  // 是否重写了重要函数

  // 处理注册页面
  extend(SignUpModal.prototype, 'submitData', function (data) {
    if (recaptcha) {
      data['captcha-solution'] = recaptcha.getSolution();
      data['captcha-id'] = recaptcha.getCaptchaID();
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
          data,
        }).then(
          (data) => {
            if (data.retcode !== 200) {
              recaptcha.setCaptchaID(data.newCaptchaID);
              that.alert = new Alert({
                type: 'error',
                children: data.retmsg,
              });
            } else {
              that.alert = new Alert({
                type: 'success',
                children: data.retmsg + "3秒后跳转到登录",
              });
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
    fields.add('recaptcha', recaptcha, -5);
  });

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

        console.log(password, remember, identification, recaptcha.getCaptchaID());
        const that = this;
        app.session
          .login({
            identification,
            'password': md5(password),
            remember,
            'captcha-solution': recaptcha.getSolution(),
            'captcha-id': recaptcha.getCaptchaID(),
          }, { errorHandler: this.onerror.bind(this) })
          .then((data) => {
            if (data.retcode !== 200) {
              recaptcha.setCaptchaID(data.newCaptchaID);
              that.alert = new Alert({
                type: 'error',
                children: data.retmsg,
              });
            } else {
              that.alert = new Alert({
                type: 'success',
                children: data.retmsg + "3秒后刷新页面",
              });
              setTimeout(() => {
                window.location.reload();
              }, 3000);
            }
            this.loading = false;
            m.redraw();
          });
      });
    }

    fields.add('recaptcha', recaptcha, -5);
  })

});