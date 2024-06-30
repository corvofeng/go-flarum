import Component from 'flarum/Component';
import Stream from 'mithril/stream';
import Captcha from "../captcha";

// 如果想使用Google的验证码, 参考这个代码
// https://github.com/FriendsOfFlarum/recaptcha/blob/master/js/src/forum/components/Recaptcha.js
export default class SelfCaptcha extends Component {
  oninit(vnode) {
    super.oninit(vnode);
    this.isLoadCaptcha = false;
    this.captchaId = "";
    console.log();

    vnode.attrs.refreshBinder(this.refreshCaptcha.bind(this));
    this.refreshCaptcha();
  }

  view() {
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
          type="text" placeholder="验证码"
          oninput={(e) => { Captcha.setCaptchaSoultion(e.target.value) }}
        />
      </div>
    );
  }

  getUrl(cid) {
    if (cid === "") return "";
    return `/static/captcha/${cid}.png`
  }
  setCaptchaID(cid) {
    this.captchaId = cid;
    this.isLoadCaptcha = true;
    Captcha.setCaptchaId(cid);
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


