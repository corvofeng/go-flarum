var Captcha = {
  captchaId: "",
  captchaSolution: "",
  refreshCallback: null,
  setCaptchaId: function (value) {
    Captcha.captchaId = value;
  },
  setCaptchaSoultion: function (value) {
    Captcha.captchaSolution = value;
  },
  bindRefreshCallback(callback) {
    Captcha.refreshCallback = callback;
  },
  onrefresh: function() {
    if (Captcha.refreshCallback) {
      Captcha.refreshCallback();
    }
  }
}

// module.exports = Captcha;
export default Captcha;