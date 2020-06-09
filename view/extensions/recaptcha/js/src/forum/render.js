// 此文件来自于线上的dist代码中的s9e, 没什么问题暂时不改了
// TODO: 查看 https://github.com/s9e/TextFormatter/
(function () {

  var HINT = {};
  HINT.EMOTICONS_NOT_AFTER = 0;
  HINT.LITEDOWN_DECODE_HTML_ENTITIES = 0;
  HINT.PREG_HAS_PASSTHROUGH = !1;
  HINT.RULE_AUTO_CLOSE = 1;
  HINT.RULE_AUTO_REOPEN = 1;
  HINT.RULE_BREAK_PARAGRAPH = 1;
  HINT.RULE_CREATE_PARAGRAPHS = 1;
  HINT.RULE_DISABLE_AUTO_BR = 1;
  HINT.RULE_ENABLE_AUTO_BR = 1;
  HINT.RULE_IGNORE_TAGS = 1;
  HINT.RULE_IGNORE_TEXT = 1;
  HINT.RULE_IGNORE_WHITESPACE = 1;
  HINT.RULE_IS_TRANSPARENT = 1;
  HINT.RULE_PREVENT_BR = 1;
  HINT.RULE_SUSPEND_AUTO_BR = 1;
  HINT.RULE_TRIM_FIRST_LINE = 1;
  HINT.attributeDefaultValue = 0;
  HINT.closeAncestor = 0;
  HINT.closeParent = 1;
  HINT.createChild = 1;
  HINT.fosterParent = 1;
  HINT.hash = 1;
  HINT.ignoreAttrs = 1;
  HINT.namespaces = 0;
  HINT.onRender = 1;
  HINT.onUpdate = 1;
  HINT.regexp = 1;
  HINT.regexpLimit = 1;
  HINT.requireAncestor = 0;
  var o553403F9 = [0];
  var oC06C5BF5 = [8481];
  var oD363F9C5 = [32551];
  var oF6AF222C = [32623];
  var o3C789569 = [15661];
  var oB565876D = {
    flags: 2
  };
  var o939F1698 = {
    flags: 514
  };
  var oA80287CC = {
    flags: 3089
  };
  var o1BC3EAF4 = {
    filterChain: [],
    required: !1
  };
  var o02D8DBB5 = {
    filterChain: [],
    required: !0
  };
  var oD4869BFF = {
    "B": 1,
    "C": 1,
    "COLOR": 1,
    "EM": 1,
    "EMAIL": 1,
    "I": 1,
    "S": 1,
    "SIZE": 1,
    "STRONG": 1,
    "U": 1,
    "URL": 1
  };
  var oDF43256F = [function (attrValue, attrName) {
    return NumericFilter.filterUint(attrValue)
  }];
  var o118B31AC = function (tag, tagConfig) {
    return filterAttributes(tag, tagConfig, registeredVars, logger)
  };
  var o1B4F5B52 = [o118B31AC];
  var o75AB9259 = {
    filterChain: oDF43256F,
    required: !1
  };
  var o6CB0A318 = {
    filterChain: oDF43256F,
    required: !0
  };
  var oF307D35C = {
    filterChain: [function (attrValue, attrName) {
      return UrlFilter.filter(attrValue, registeredVars.urlConfig, logger)
    }],
    required: !0
  };
  var o1AF69066 = {
    closeParent: oD4869BFF,
    flags: 268,
    fosterParent: oD4869BFF
  };
  var o8F52476D = {
    allowed: oD363F9C5,
    attributes: {},
    bitNumber: 2,
    filterChain: o1B4F5B52,
    nestingLimit: 10,
    rules: {
      flags: 0
    },
    tagLimit: 5000
  };
  var o5B6ED7AA = {
    allowed: oD363F9C5,
    attributes: {},
    bitNumber: 2,
    filterChain: o1B4F5B52,
    nestingLimit: 10,
    rules: oB565876D,
    tagLimit: 5000
  };
  var oE13673F5 = {
    allowed: oD363F9C5,
    attributes: {},
    bitNumber: 3,
    filterChain: o1B4F5B52,
    nestingLimit: 10,
    rules: {
      closeParent: oD4869BFF,
      flags: 260,
      fosterParent: oD4869BFF
    },
    tagLimit: 5000
  };
  var xsl = (
    "<xsl:stylesheet version=\"2.0\" xmlns:xsl=\"http:\/\/www.w3.org\/1999\/XSL\/Transform\"><xsl:output method=\"html\" encoding=\"utf-8\" indent=\"no\"\/><xsl:param$kDISCUSSION_URL\">https:\/\/flarum.yjzq.fun\/d\/<\/xsl:param><xsl:param$kPROFILE_URL\">https:\/\/flarum.yjzq.fun\/u\/<\/xsl:param>$aB|DEL|EM|H1|H2|H3|H4|H5|H6|I|LI|S|STRONG|SUB|SUP|U|p\"><xsl:element$k{translate(name(),'BDEGHILMNOPRSTU','bdeghilmnoprstu')}\">$b<\/xsl:element>$c$aC\"><code>$b<\/code>$c$aCENTER\"><div$itext-align:center\">$b<\/div>$c$aCODE\"><pre data-s9e-livepreview-hash=\"\" data-s9e-livepreview-onupdate=\"if(typeof hljsLoader!=='undefined')hljsLoader.highlightBlocks(this)\"><code><xsl:if$e@lang\"><xsl:attribute$kclass\">language-$d@lang\"\/><\/xsl:attribute><\/xsl:if>$b<\/code><script async=\"\" crossorigin=\"anonymous\" data-hljs-style=\"github-gist\" data-s9e-livepreview-onrender=\"if(typeof hljsLoader!=='undefined')this.parentNode.removeChild(this)\" integrity=\"sha384-PG1zopchh98J\/8eUG\/5ESEA+dE1Art6Ym1XKxLljmoOwlodOfUguCC\/cydPWZSQ\/\" onload=\"hljsLoader.highlightBlocks(this.parentNode)\" src=\"https:\/\/cdn.jsdelivr.net\/gh\/s9e\/hljs-loader@1.0.8\/loader.min.js\"\/><\/pre>$c$aCOLOR\"><span$icolor:{@color}\">$b<\/span>$c$aE\"><xsl:choose><$f$e.=':)'\">\ud83d\ude42<\/$f><$f$e.=':D'\">\ud83d\ude03<\/$f><$f$e.=':P'\">\ud83d\ude1b<\/$f><$f$e.=':('\">\ud83d\ude41<\/$f><$f$e.=':|'\">\ud83d\ude10<\/$f><$f$e.=';)'\">\ud83d\ude09<\/$f><$f$e.=&quot;:'(&quot;\">\ud83d\ude22<\/$f><$f$e.=':O'\">\ud83d\ude2e<\/$f><$f$e.='&gt;:('\">\ud83d\ude21<\/$f><xsl:otherwise>$d.\"\/><\/xsl:otherwise><\/xsl:choose>$c$aEMAIL\"><a href=\"mailto:{@email}\">$b<\/a>$c$aESC\">$b$c$aHR\"><hr\/>$c$aIMG\"><img src=\"{@src}\" title=\"{@title}\" alt=\"{@alt}\">$gheight|@width\"\/><\/img>$c$aISPOILER\"><span$lspoiler\"$hclass\" onclick=\"removeAttribute('class')\">$b<\/span>$c$aLIST\"><xsl:choose><$f$enot(@type)\"><ul>$b<\/ul><\/$f><$f$estarts-with(@type,'decimal')or starts-with(@type,'lower')or starts-with(@type,'upper')\"><ol$i$j>$gstart\"\/>$b<\/ol><\/$f><xsl:otherwise><ul$i$j>$b<\/ul><\/xsl:otherwise><\/xsl:choose>$c$aPOSTMENTION\"><a href=\"{$DISCUSSION_URL}{@discussionid}\/{@number}\"$lPostMention\" data-id=\"{@id}\">$d@dis" + "playname\"\/><\/a>$c$aQUOTE\"><blockquote><xsl:if$enot(@author)\"><xsl:attribute$kclass\">uncited<\/xsl:attribute><\/xsl:if><div><xsl:if$e@author\"><cite>$d@author\"\/> wrote:<\/cite><\/xsl:if>$b<\/div><\/blockquote>$c$aSIZE\"><span$ifont-size:{@size}px\">$b<\/span>$c$aSPOILER\"><details$lspoiler\"$hopen\">$b<\/details>$c$aURL\"><a href=\"{@url}\" rel=\" nofollow ugc\">$gtitle\"\/>$b<\/a>$c$aUSERMENTION\"><a href=\"{$PROFILE_URL}{@username}\"$lUserMention\">@$d@displayname\"\/><\/a>$c$abr\"><br\/>$c$ae|i|s\"\/><\/xsl:stylesheet>").replace(/\$[a-l]/g, function (k) {
    return {
      "$a": "<xsl:template match=\"",
      "$b": "<xsl:apply-templates\/>",
      "$c": "<\/xsl:template>",
      "$d": "<xsl:value-of select=\"",
      "$e": " test=\"",
      "$f": "xsl:when",
      "$g": "<xsl:copy-of select=\"@",
      "$h": " data-s9e-livepreview-ignore-attrs=\"",
      "$i": " style=\"",
      "$j": "list-style-type:{@type}\"",
      "$k": " name=\"",
      "$l": " class=\""
    } [k]
  });
  console.log(xsl);

  var EmailFilter = {
    filter: function (attrValue) {
      return /^[-\w.+]+@[-\w.]+$/.test(attrValue) ? attrValue : !1
    }
  };
  var FalseFilter = {
    filter: function (attrValue) {
      return !1
    }
  };
  var HashmapFilter = {
    filter: function (attrValue, map, strict) {
      if (attrValue in map) {
        return map[attrValue]
      }
      return (strict) ? !1 : attrValue
    }
  };
  var MapFilter = {
    filter: function (attrValue, map) {
      var i = -1,
        cnt = map.length;
      while (++i < cnt) {
        if (map[i][0].test(attrValue)) {
          return map[i][1]
        }
      }
      return attrValue
    }
  };
  var NetworkFilter = {
    filterIp: function (attrValue) {
      if (/^[\d.]+$/.test(attrValue)) {
        return NetworkFilter.filterIpv4(attrValue)
      }
      if (/^[\da-f:]+$/i.test(attrValue)) {
        return NetworkFilter.filterIpv6(attrValue)
      }
      return !1
    },
    filterIpport: function (attrValue) {
      var m, ip;
      if (m = /^\[([\da-f:]+)(\]:[1-9]\d*)$/i.exec(attrValue)) {
        ip = NetworkFilter.filterIpv6(m[1]);
        if (ip === !1) {
          return !1
        }
        return '[' + ip + m[2]
      }
      if (m = /^([\d.]+)(:[1-9]\d*)$/.exec(attrValue)) {
        ip = NetworkFilter.filterIpv4(m[1]);
        if (ip === !1) {
          return !1
        }
        return ip + m[2]
      }
      return !1
    },
    filterIpv4: function (attrValue) {
      if (!/^\d+\.\d+\.\d+\.\d+$/.test(attrValue)) {
        return !1
      }
      var i = 4,
        p = attrValue.split('.');
      while (--i >= 0) {
        if (p[i][0] === '0' || p[i] > 255) {
          return !1
        }
      }
      return attrValue
    },
    filterIpv6: function (attrValue) {
      return /^([\da-f]{0,4}:){2,7}(?:[\da-f]{0,4}|\d+\.\d+\.\d+\.\d+)$/.test(attrValue) ? attrValue : !1
    }
  };
  var NumericFilter = {
    filterFloat: function (attrValue) {
      return /^(?:0|-?[1-9]\d*)(?:\.\d+)?(?:e[1-9]\d*)?$/i.test(attrValue) ? attrValue : !1
    },
    filterInt: function (attrValue) {
      return /^(?:0|-?[1-9]\d*)$/.test(attrValue) ? attrValue : !1
    },
    filterRange: function (attrValue, min, max, logger) {
      if (!/^(?:0|-?[1-9]\d*)$/.test(attrValue)) {
        return !1
      }
      attrValue = parseInt(attrValue, 10);
      if (attrValue < min) {
        if (logger) {
          logger.warn('Value outside of range, adjusted up to min value', {
            'attrValue': attrValue,
            'min': min,
            'max': max
          })
        }
        return min
      }
      if (attrValue > max) {
        if (logger) {
          logger.warn('Value outside of range, adjusted down to max value', {
            'attrValue': attrValue,
            'min': min,
            'max': max
          })
        }
        return max
      }
      return attrValue
    },
    filterUint: function (attrValue) {
      return /^(?:0|[1-9]\d*)$/.test(attrValue) ? attrValue : !1
    }
  };
  var RegexpFilter = {
    filter: function (attrValue, regexp) {
      return regexp.test(attrValue) ? attrValue : !1
    }
  };
  var TimestampFilter = {
    filter: function (attrValue) {
      var m = /^(?=\d)(?:(\d+)h)?(?:(\d+)m)?(?:(\d+)s)?$/.exec(attrValue);
      if (m) {
        return 3600 * (m[1] || 0) + 60 * (m[2] || 0) + (+m[3] || 0)
      }
      return NumericFilter.filterUint(attrValue)
    }
  };
  var UrlFilter = {
    filter: function (attrValue, urlConfig, logger) {
      var p = UrlFilter.parseUrl(attrValue.replace(/^\s+/, '').replace(/\s+$/, ''));
      var error = UrlFilter.validateUrl(urlConfig, p);
      if (error) {
        if (logger) {
          p.attrValue = attrValue;
          logger.err(error, p)
        }
        return !1
      }
      return UrlFilter.rebuildUrl(urlConfig, p)
    },
    parseUrl: function (url) {
      var regexp = /^(?:([a-z][-+.\w]*):)?(?:\/\/(?:([^:\/?#]*)(?::([^\/?#]*)?)?@)?(?:(\[[a-f\d:]+\]|[^:\/?#]+)(?::(\d*))?)?(?![^\/?#]))?([^?#]*)(\?[^#]*)?(#.*)?$/i;
      var m = regexp.exec(url),
        parts = {},
        tokens = ['scheme', 'user', 'pass', 'host', 'port', 'path', 'query', 'fragment'];
      tokens.forEach(function (name, i) {
        parts[name] = (m[i + 1] > '') ? m[i + 1] : ''
      });
      parts.scheme = parts.scheme.toLowerCase();
      parts.host = parts.host.replace(/[\u3002\uff0e\uff61]/g, '.').replace(/\.+$/g, '');
      if (/[^\x00-\x7F]/.test(parts.host) && typeof punycode !== 'undefined') {
        parts.host = punycode.toASCII(parts.host)
      }
      return parts
    },
    rebuildUrl: function (urlConfig, p) {
      var url = '';
      if (p.scheme !== '') {
        url += p.scheme + ':'
      }
      if (p.host !== '') {
        url += '//';
        if (p.user !== '') {
          url += rawurlencode(decodeURIComponent(p.user));
          if (p.pass !== '') {
            url += ':' + rawurlencode(decodeURIComponent(p.pass))
          }
          url += '@'
        }
        url += p.host;
        if (p.port !== '') {
          url += ':' + p.port
        }
      } else if (p.scheme === 'file') {
        url += '//'
      }
      var path = p.path + p.query + p.fragment;
      path = path.replace(/%.?[a-f]/g, function (str) {
        return str.toUpperCase()
      }, path);
      url += UrlFilter.sanitizeUrl(path);
      if (!p.scheme) {
        url = url.replace(/^([^\/]*):/, '$1%3A')
      }
      return url
    },
    sanitizeUrl: function (url) {
      return url.replace(/[^\u0020-\u007E]+/g, encodeURIComponent).replace(/%(?![0-9A-Fa-f]{2})|[^!#-&*-;=?-Z_a-z~]/g, escape)
    },
    validateUrl: function (urlConfig, p) {
      if (p.scheme !== '' && !urlConfig.allowedSchemes.test(p.scheme)) {
        return 'URL scheme is not allowed'
      }
      if (p.host !== '') {
        var regexp = /^(?!-)[-a-z0-9]{0,62}[a-z0-9](?:\.(?!-)[-a-z0-9]{0,62}[a-z0-9])*$/i;
        if (!regexp.test(p.host)) {
          if (!NetworkFilter.filterIpv4(p.host) && !NetworkFilter.filterIpv6(p.host.replace(/^\[(.*)\]$/, '$1', p.host))) {
            return 'URL host is invalid'
          }
        }
        if ((urlConfig.disallowedHosts && urlConfig.disallowedHosts.test(p.host)) || (urlConfig.restrictedHosts && !urlConfig.restrictedHosts.test(p.host))) {
          return 'URL host is not allowed'
        }
      } else if (/^(?:(?:f|ht)tps?)$/.test(p.scheme)) {
        return 'Missing host'
      }
    }
  };

  function html_entity_decode(str) {
    var b = document.createElement('b');
    html_entity_decode = function (str) {
      b.innerHTML = str.replace(/</g, '&lt;');
      return b.textContent
    };
    return html_entity_decode(str)
  }

  function htmlspecialchars_compat(str) {
    var t = {
      '<': '&lt;',
      '>': '&gt;',
      '&': '&amp;',
      '"': '&quot;'
    };
    return str.replace(/[<>&"]/g, function (c) {
      return t[c]
    })
  }

  function htmlspecialchars_noquotes(str) {
    var t = {
      '<': '&lt;',
      '>': '&gt;',
      '&': '&amp;'
    };
    return str.replace(/[<>&]/g, function (c) {
      return t[c]
    })
  }

  function rawurlencode(str) {
    return encodeURIComponent(str).replace(/[!'()*]/g, function (c) {
      return '%' + c.charCodeAt(0).toString(16).toUpperCase()
    })
  }

  function returnFalse() {
    return !1
  }

  function returnTrue() {
    return !0
  }

  function executeAttributePreprocessors(tag, tagConfig) {
    if (typeof tagConfig.attributePreprocessors === 'undefined') {
      return
    }
    tagConfig.attributePreprocessors.forEach(function (ap) {
      var attrName = ap[0],
        regexp = ap[1],
        map = ap[2];
      if (tag.hasAttribute(attrName)) {
        executeAttributePreprocessor(tag, attrName, regexp, map)
      }
    })
  }

  function filterAttributes(tag, tagConfig, registeredVars, logger) {
    var attributes = {},
      attrName;
    for (attrName in tagConfig.attributes) {
      var attrConfig = tagConfig.attributes[attrName],
        attrValue = !1;
      if (tag.hasAttribute(attrName)) {
        attrValue = executeAttributeFilterChain(attrConfig.filterChain, attrName, tag.getAttribute(attrName))
      }
      if (attrValue !== !1) {
        attributes[attrName] = attrValue
      } else if (HINT.attributeDefaultValue && typeof attrConfig.defaultValue !== 'undefined') {
        attributes[attrName] = attrConfig.defaultValue
      } else if (attrConfig.required) {
        tag.invalidate()
      }
    }
    tag.setAttributes(attributes)
  }

  function filterTag(tag) {
    var tagName = tag.getName(),
      tagConfig = tagsConfig[tagName];
    logger.setTag(tag);
    for (var i = 0; i < tagConfig.filterChain.length; ++i) {
      if (tag.isInvalid()) {
        break
      }
      tagConfig.filterChain[i](tag, tagConfig)
    }
    logger.unsetTag()
  }

  function executeAttributeFilterChain(filterChain, attrName, attrValue) {
    logger.setAttribute(attrName);
    for (var i = 0; i < filterChain.length; ++i) {
      attrValue = filterChain[i](attrValue, attrName);
      if (attrValue === !1) {
        break
      }
    }
    logger.unsetAttribute();
    return attrValue
  }

  function executeAttributePreprocessor(tag, attrName, regexp, map) {
    var attrValue = tag.getAttribute(attrName),
      captures = getNamedCaptures(attrValue, regexp, map),
      k;
    for (k in captures) {
      if (k === attrName || !tag.hasAttribute(k)) {
        tag.setAttribute(k, captures[k])
      }
    }
  }

  function getNamedCaptures(attrValue, regexp, map) {
    var m = regexp.exec(attrValue);
    if (!m) {
      return []
    }
    var values = {};
    map.forEach(function (k, i) {
      if (typeof m[i] === 'string' && m[i] !== '') {
        values[k] = m[i]
      }
    });
    return values
  }

  function Logger() {}
  Logger.prototype = {
    add: function (type, msg, context) {},
    clear: function () {},
    setAttribute: function (attrName) {},
    setTag: function (tag) {},
    unsetAttribute: function () {},
    unsetTag: function () {},
    debug: function (msg, context) {},
    err: function (msg, context) {},
    info: function (msg, context) {},
    warn: function (msg, context) {}
  };

  function Tag(type, name, pos, len, priority) {
    this.type = +type;
    this.name = name;
    this.pos = +pos;
    this.len = +len;
    this.sortPriority = +priority || 0;
    this.attributes = {};
    this.cascade = [];
    if (isNaN(type + pos + len)) {
      this.invalidate()
    }
  }
  Tag.START_TAG = 1;
  Tag.END_TAG = 2;
  Tag.SELF_CLOSING_TAG = 3;
  Tag.prototype.attributes;
  Tag.prototype.cascade;
  Tag.prototype.endTag;
  Tag.prototype.invalid = !1;
  Tag.prototype.len;
  Tag.prototype.name;
  Tag.prototype.pos;
  Tag.prototype.sortPriority;
  Tag.prototype.startTag;
  Tag.prototype.type;
  Tag.prototype.addFlags = function (flags) {
    this.flags |= flags
  };
  Tag.prototype.cascadeInvalidationTo = function (tag) {
    this.cascade.push(tag);
    if (this.invalid) {
      tag.invalidate()
    }
  };
  Tag.prototype.invalidate = function () {
    if (!this.invalid) {
      this.invalid = !0;
      this.cascade.forEach(function (tag) {
        tag.invalidate()
      })
    }
  };
  Tag.prototype.pairWith = function (tag) {
    if (this.canBePaired(this, tag)) {
      this.endTag = tag;
      tag.startTag = this;
      this.cascadeInvalidationTo(tag)
    } else if (this.canBePaired(tag, this)) {
      this.startTag = tag;
      tag.endTag = this
    }
  };
  Tag.prototype.canBePaired = function (startTag, endTag) {
    return startTag.name === endTag.name && startTag.type === Tag.START_TAG && endTag.type === Tag.END_TAG && startTag.pos <= startTag.pos
  };
  Tag.prototype.removeFlags = function (flags) {
    this.flags &= ~flags
  };
  Tag.prototype.setFlags = function (flags) {
    this.flags = flags
  };
  Tag.prototype.getAttributes = function () {
    var attributes = {};
    for (var attrName in this.attributes) {
      attributes[attrName] = this.attributes[attrName]
    }
    return attributes
  };
  Tag.prototype.getEndTag = function () {
    return this.endTag
  };
  Tag.prototype.getFlags = function () {
    return this.flags
  };
  Tag.prototype.getLen = function () {
    return this.len
  };
  Tag.prototype.getName = function () {
    return this.name
  };
  Tag.prototype.getPos = function () {
    return this.pos
  };
  Tag.prototype.getSortPriority = function () {
    return this.sortPriority
  };
  Tag.prototype.getStartTag = function () {
    return this.startTag
  };
  Tag.prototype.getType = function () {
    return this.type
  };
  Tag.prototype.canClose = function (startTag) {
    if (this.invalid || !this.canBePaired(startTag, this) || (this.startTag && this.startTag !== startTag) || (startTag.endTag && startTag.endTag !== this)) {
      return !1
    }
    return !0
  };
  Tag.prototype.isBrTag = function () {
    return (this.name === 'br')
  };
  Tag.prototype.isEndTag = function () {
    return !!(this.type & Tag.END_TAG)
  };
  Tag.prototype.isIgnoreTag = function () {
    return (this.name === 'i')
  };
  Tag.prototype.isInvalid = function () {
    return this.invalid
  };
  Tag.prototype.isParagraphBreak = function () {
    return (this.name === 'pb')
  };
  Tag.prototype.isSelfClosingTag = function () {
    return (this.type === Tag.SELF_CLOSING_TAG)
  };
  Tag.prototype.isSystemTag = function () {
    return ('br i pb v'.indexOf(this.name) > -1)
  };
  Tag.prototype.isStartTag = function () {
    return !!(this.type & Tag.START_TAG)
  };
  Tag.prototype.isVerbatim = function () {
    return (this.name === 'v')
  };
  Tag.prototype.getAttribute = function (attrName) {
    return this.attributes[attrName]
  };
  Tag.prototype.hasAttribute = function (attrName) {
    return (attrName in this.attributes)
  };
  Tag.prototype.removeAttribute = function (attrName) {
    delete this.attributes[attrName]
  };
  Tag.prototype.setAttribute = function (attrName, attrValue) {
    this.attributes[attrName] = attrValue
  };
  Tag.prototype.setAttributes = function (attributes) {
    this.attributes = {};
    for (var attrName in attributes) {
      this.attributes[attrName] = attributes[attrName]
    }
  };
  var RULE_AUTO_CLOSE = 1 << 0;
  var RULE_AUTO_REOPEN = 1 << 1;
  var RULE_BREAK_PARAGRAPH = 1 << 2;
  var RULE_CREATE_PARAGRAPHS = 1 << 3;
  var RULE_DISABLE_AUTO_BR = 1 << 4;
  var RULE_ENABLE_AUTO_BR = 1 << 5;
  var RULE_IGNORE_TAGS = 1 << 6;
  var RULE_IGNORE_TEXT = 1 << 7;
  var RULE_IGNORE_WHITESPACE = 1 << 8;
  var RULE_IS_TRANSPARENT = 1 << 9;
  var RULE_PREVENT_BR = 1 << 10;
  var RULE_SUSPEND_AUTO_BR = 1 << 11;
  var RULE_TRIM_FIRST_LINE = 1 << 12;
  var RULES_AUTO_LINEBREAKS = RULE_DISABLE_AUTO_BR | RULE_ENABLE_AUTO_BR | RULE_SUSPEND_AUTO_BR;
  var RULES_INHERITANCE = RULE_ENABLE_AUTO_BR;
  var WHITESPACE = " \n\t";
  var cntOpen;
  var cntTotal;
  var context;
  var currentFixingCost;
  var currentTag;
  var isRich;
  var logger = new Logger;
  var maxFixingCost = 10000;
  var namespaces;
  var openTags;
  var output;
  var plugins = {
    "Autoemail": {
      parser: function (text, matches) {
        var config = {
          attrName: "email",
          tagName: "EMAIL"
        };
        var tagName = config.tagName,
          attrName = config.attrName;
        matches.forEach(function (m) {
          var startTag = addStartTag(tagName, m[0][1], 0);
          startTag.setAttribute(attrName, m[0][0]);
          var endTag = addEndTag(tagName, m[0][1] + m[0][0].length, 0);
          startTag.pairWith(endTag)
        })
      },
      quickMatch: "@",
      regexp: /\b[-a-z0-9_+.]+@[-a-z0-9.]*[a-z0-9]/ig,
      regexpLimit: 50000
    },
    "Autolink": {
      parser: function (text, matches) {
        var config = {
          attrName: "url",
          tagName: "URL"
        };
        matches.forEach(function (m) {
          linkifyUrl(m[0][1], trimUrl(m[0][0]))
        });

        function linkifyUrl(tagPos, url) {
          if (!/^www\.|^[^:]+:/i.test(url)) {
            return
          }
          var endPos = tagPos + url.length,
            endTag = addEndTag(config.tagName, endPos, 0);
          if (url[3] === '.') {
            url = 'http://' + url
          }
          var startTag = addStartTag(config.tagName, tagPos, 0, 1);
          startTag.setAttribute(config.attrName, url);
          startTag.pairWith(endTag);
          var contentTag = addVerbatim(tagPos, endPos - tagPos, 1000);
          startTag.cascadeInvalidationTo(contentTag)
        }

        function trimUrl(url) {
          return url.replace(/(?:(?![-=)\/_])[\s!-.:-@[-`{-~])+$/, '')
        }
      },
      quickMatch: ":",
      regexp: /\bhttps?:(?:[^\s()\[\]\uFF01-\uFF0F\uFF1A-\uFF20\uFF3B-\uFF40\uFF5B-\uFF65]|\([^\s()]*\)|\[\w*\])+/ig,
      regexpLimit: 50000
    },
    "BBCodes": {
      parser: function (text, matches) {
        var config = {
          bbcodes: {
            "*": {
              tagName: "LI"
            },
            "B": [],
            "CENTER": [],
            "CODE": {
              defaultAttribute: "lang"
            },
            "COLOR": [],
            "DEL": [],
            "EMAIL": {
              contentAttributes: ["email"]
            },
            "I": [],
            "IMG": {
              contentAttributes: ["src"],
              defaultAttribute: "src"
            },
            "LIST": {
              defaultAttribute: "type"
            },
            "QUOTE": {
              defaultAttribute: "author"
            },
            "S": [],
            "SIZE": [],
            "U": [],
            "URL": {
              contentAttributes: ["url"]
            }
          }
        };
        var attributes;
        var bbcodeConfig;
        var bbcodeName;
        var bbcodeSuffix;
        var pos;
        var startPos;
        var textLen = text.length;
        var uppercaseText = '';
        matches.forEach(function (m) {
          bbcodeName = m[1][0].toUpperCase();
          if (!(bbcodeName in config.bbcodes)) {
            return
          }
          bbcodeConfig = config.bbcodes[bbcodeName];
          startPos = m[0][1];
          pos = startPos + m[0][0].length;
          try {
            parseBBCode()
          } catch (e) {}
        });

        function addBBCodeEndTag() {
          return addEndTag(getTagName(), startPos, pos - startPos)
        }

        function addBBCodeSelfClosingTag() {
          var tag = addSelfClosingTag(getTagName(), startPos, pos - startPos);
          tag.setAttributes(attributes);
          return tag
        }

        function addBBCodeStartTag() {
          var tag = addStartTag(getTagName(), startPos, pos - startPos);
          tag.setAttributes(attributes);
          return tag
        }

        function captureEndTag() {
          if (!uppercaseText) {
            uppercaseText = text.toUpperCase()
          }
          var match = '[/' + bbcodeName + bbcodeSuffix + ']',
            endTagPos = uppercaseText.indexOf(match, pos);
          if (endTagPos < 0) {
            return null
          }
          return addEndTag(getTagName(), endTagPos, match.length)
        }

        function getTagName() {
          return bbcodeConfig.tagName || bbcodeName
        }

        function parseAttributes() {
          var firstPos = pos,
            attrName;
          attributes = {};
          while (pos < textLen) {
            var c = text[pos];
            if (" \n\t".indexOf(c) > -1) {
              ++pos;
              continue
            }
            if ('/]'.indexOf(c) > -1) {
              return
            }
            var spn = /^[-\w]*/.exec(text.substr(pos, 100))[0].length;
            if (spn) {
              attrName = text.substr(pos, spn).toLowerCase();
              pos += spn;
              if (pos >= textLen) {
                throw ''
              }
              if (text[pos] !== '=') {
                continue
              }
            } else if (c === '=' && pos === firstPos) {
              attrName = bbcodeConfig.defaultAttribute || bbcodeName.toLowerCase()
            } else {
              throw ''
            }
            if (++pos >= textLen) {
              throw ''
            }
            attributes[attrName] = parseAttributeValue()
          }
        }

        function parseAttributeValue() {
          if (text[pos] === '"' || text[pos] === "'") {
            return parseQuotedAttributeValue()
          }
          var match = /(?:[^\s\]]|[ \t](?!\s*(?:[-\w]+=|\/?\])))*/.exec(text.substr(pos)),
            attrValue = match[0];
          pos += attrValue.length;
          return attrValue
        }

        function parseBBCode() {
          parseBBCodeSuffix();
          if (text[startPos + 1] === '/') {
            if (text[pos] === ']' && bbcodeSuffix === '') {
              ++pos;
              addBBCodeEndTag()
            }
            return
          }
          parseAttributes();
          if (bbcodeConfig.predefinedAttributes) {
            for (var attrName in bbcodeConfig.predefinedAttributes) {
              if (!(attrName in attributes)) {
                attributes[attrName] = bbcodeConfig.predefinedAttributes[attrName]
              }
            }
          }
          if (text[pos] === ']') {
            ++pos
          } else {
            if (text.substr(pos, 2) === '/]') {
              pos += 2;
              addBBCodeSelfClosingTag()
            }
            return
          }
          var contentAttributes = [];
          if (bbcodeConfig.contentAttributes) {
            bbcodeConfig.contentAttributes.forEach(function (attrName) {
              if (!(attrName in attributes)) {
                contentAttributes.push(attrName)
              }
            })
          }
          var requireEndTag = (bbcodeSuffix || bbcodeConfig.forceLookahead),
            endTag = (requireEndTag || contentAttributes.length) ? captureEndTag() : null;
          if (endTag) {
            contentAttributes.forEach(function (attrName) {
              attributes[attrName] = text.substr(pos, endTag.getPos() - pos)
            })
          } else if (requireEndTag) {
            return
          }
          var tag = addBBCodeStartTag();
          if (endTag) {
            tag.pairWith(endTag)
          }
        }

        function parseBBCodeSuffix() {
          bbcodeSuffix = '';
          if (text[pos] === ':') {
            bbcodeSuffix = /^:\d*/.exec(text.substr(pos))[0];
            pos += bbcodeSuffix.length
          }
        }

        function parseQuotedAttributeValue() {
          var quote = text[pos],
            valuePos = pos + 1;
          do {
            pos = text.indexOf(quote, pos + 1);
            if (pos < 0) {
              throw ''
            }
            var n = 1;
            while (text[pos - n] === '\\') {
              ++n
            }
          }
          while (n % 2 === 0);
          var attrValue = text.substr(valuePos, pos - valuePos);
          if (attrValue.indexOf('\\') > -1) {
            attrValue = attrValue.replace(/\\([\\'"])/g, '$1')
          }
          ++pos;
          return attrValue
        }
      },
      quickMatch: "[",
      regexp: /\[\/?(\*|[-\w]+)(?=[\]\s=:\/])/g,
      regexpLimit: 50000
    },
    "Emoticons": {
      parser: function (text, matches) {
        var config = {
          tagName: "E"
        };
        matches.forEach(function (m) {
          if (HINT.EMOTICONS_NOT_AFTER && config.notAfter && m[0][1] && config.notAfter.test(text[m[0][1] - 1])) {
            return
          }
          addSelfClosingTag(config.tagName, m[0][1], m[0][0].length)
        })
      },
      regexp: /(?::(?:[()DOP|]|'\()|;\)|>:\()/g,
      regexpLimit: 50000
    },
    "Escaper": {
      parser: function (text, matches) {
        var config = {
          tagName: "ESC"
        };
        matches.forEach(function (m) {
          addTagPair(config.tagName, m[0][1], 1, m[0][1] + m[0][0].length, 0)
        })
      },
      quickMatch: "\\",
      regexp: /\\[-!#()*+.:<>@[\\\]^_`{|}~]/g,
      regexpLimit: 50000
    },
    "Litedown": {
      parser: function (text, matches) {
        var config = {
          decodeHtmlEntities: !1
        };
        var decodeHtmlEntities = config.decodeHtmlEntities;
        var hasEscapedChars = !1;
        var hasReferences = !1;
        var linkReferences = {};
        if (text.indexOf('\\') >= 0) {
          hasEscapedChars = !0;
          text = text.replace(/\\[!"'()*<>[\\\]^_`~]/g, function (str) {
            return {
              '\\!': "\x1B0",
              '\\"': "\x1B1",
              "\\'": "\x1B2",
              '\\(': "\x1B3",
              '\\)': "\x1B4",
              '\\*': "\x1B5",
              '\\<': "\x1B6",
              '\\>': "\x1B7",
              '\\[': "\x1B8",
              '\\\\': "\x1B9",
              '\\]': "\x1BA",
              '\\^': "\x1BB",
              '\\_': "\x1BC",
              '\\`': "\x1BD",
              '\\~': "\x1BE"
            } [str]
          })
        }
        text += "\n\n\x17";

        function decode(str) {
          if (HINT.LITEDOWN_DECODE_HTML_ENTITIES && decodeHtmlEntities && str.indexOf('&') > -1) {
            str = html_entity_decode(str)
          }
          str = str.replace(/\x1A/g, '');
          if (hasEscapedChars) {
            str = str.replace(/\x1B./g, function (seq) {
              return {
                "\x1B0": '!',
                "\x1B1": '"',
                "\x1B2": "'",
                "\x1B3": '(',
                "\x1B4": ')',
                "\x1B5": '*',
                "\x1B6": '<',
                "\x1B7": '>',
                "\x1B8": '[',
                "\x1B9": '\\',
                "\x1BA": ']',
                "\x1BB": '^',
                "\x1BC": '_',
                "\x1BD": '`',
                "\x1BE": '~'
              } [seq]
            })
          }
          return str
        }

        function isAfterWhitespace(pos) {
          return (pos > 0 && isWhitespace(text.charAt(pos - 1)))
        }

        function isAlnum(chr) {
          return (' abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'.indexOf(chr) > 0)
        }

        function isBeforeWhitespace(pos) {
          return isWhitespace(text[pos + 1])
        }

        function isSurroundedByAlnum(pos, len) {
          return (pos > 0 && isAlnum(text[pos - 1]) && isAlnum(text[pos + len]))
        }

        function isWhitespace(chr) {
          return (" \n\t".indexOf(chr) > -1)
        }

        function markBoundary(pos) {
          text = text.substr(0, pos) + "\x17" + text.substr(pos + 1)
        }

        function overwrite(pos, len) {
          if (len > 0) {
            text = text.substr(0, pos) + new Array(1 + len).join("\x1A") + text.substr(pos + len)
          }
        }

        function parseInlineMarkup(str, regexp, tagName) {
          if (text.indexOf(str) === -1) {
            return
          }
          var m;
          while (m = regexp.exec(text)) {
            var match = m[0],
              matchPos = m.index,
              matchLen = match.length,
              endPos = matchPos + matchLen - 2;
            addTagPair(tagName, matchPos, 2, endPos, 2);
            overwrite(matchPos, 2);
            overwrite(endPos, 2)
          }
        }

        function parseAbstractScript(tagName, syntaxChar, shortRegexp, longRegexp) {
          var pos = text.indexOf(syntaxChar);
          if (pos === -1) {
            return
          }
          parseShortForm(pos);
          parseLongForm(pos);

          function parseLongForm(pos) {
            pos = text.indexOf(syntaxChar + '(', pos);
            if (pos === -1) {
              return
            }
            var m, regexp = longRegexp;
            regexp.lastIndex = pos;
            while (m = regexp.exec(text)) {
              var match = m[0],
                matchPos = m.index,
                matchLen = match.length;
              addTagPair(tagName, matchPos, 2, matchPos + matchLen - 1, 1);
              overwrite(matchPos, matchLen)
            }
            if (match) {
              parseLongForm(pos)
            }
          }

          function parseShortForm(pos) {
            var m, regexp = shortRegexp;
            regexp.lastIndex = pos;
            while (m = regexp.exec(text)) {
              var match = m[0],
                matchPos = m.index,
                matchLen = match.length,
                startPos = matchPos,
                endLen = (match.substr(-1) === syntaxChar) ? 1 : 0,
                endPos = matchPos + matchLen - endLen;
              addTagPair(tagName, startPos, 1, endPos, endLen)
            }
          }
        }

        function setLinkAttributes(tag, linkInfo, attrName) {
          var url = linkInfo.replace(/^\s*/, '').replace(/\s*$/, ''),
            title = '',
            pos = url.indexOf(' ');
          if (pos !== -1) {
            title = url.substr(pos).replace(/^\s*\S/, '').replace(/\S\s*$/, '');
            url = url.substr(0, pos)
          }
          if (/^<.+>$/.test(url)) {
            url = url.replace(/^<(.+)>$/, '$1').replace(/\\>/g, '>')
          }
          tag.setAttribute(attrName, decode(url));
          if (title > '') {
            tag.setAttribute('title', decode(title))
          }
        }(function () {
          var setextLines = {};

          function parse() {
            matchSetextLines();
            var blocks = [],
              blocksCnt = 0,
              codeFence, codeIndent = 4,
              codeTag, lineIsEmpty = !0,
              lists = [],
              listsCnt = 0,
              newContext = !1,
              textBoundary = 0,
              breakParagraph, continuation, ignoreLen, indentStr, indentLen, lfPos, listIndex, maxIndent, minIndent, blockDepth, tagPos, tagLen;
            var matches = [],
              m, regexp = /^(?:(?=[-*+\d \t>`~#_])((?: {0,3}>(?:(?!!)|!(?![^\n>]*?!<)) ?)+)?([ \t]+)?(\* *\* *\*[* ]*$|- *- *-[- ]*$|_ *_ *_[_ ]*$)?((?:[-*+]|\d+\.)[ \t]+(?=\S))?[ \t]*(#{1,6}[ \t]+|```+[^`\n]*$|~~~+[^~\n]*$)?)?/gm;
            while (m = regexp.exec(text)) {
              matches.push(m);
              if (m.index === regexp.lastIndex) {
                ++regexp.lastIndex
              }
            }
            matches.forEach(function (m) {
              var blockMarks = [],
                matchPos = m.index,
                matchLen = m[0].length,
                startPos, startLen, endPos, endLen;
              ignoreLen = 0;
              blockDepth = 0;
              continuation = !lineIsEmpty;
              lfPos = text.indexOf("\n", matchPos);
              lineIsEmpty = (lfPos === matchPos + matchLen && !m[3] && !m[4] && !m[5]);
              if (!matchLen) {
                ++regexp.lastIndex
              }
              breakParagraph = (lineIsEmpty && continuation);
              if (m[1]) {
                blockMarks = getBlockMarks(m[1]);
                blockDepth = blockMarks.length;
                ignoreLen = m[1].length;
                if (codeTag && codeTag.hasAttribute('blockDepth')) {
                  blockDepth = Math.min(blockDepth, codeTag.getAttribute('blockDepth'));
                  ignoreLen = computeBlockIgnoreLen(m[1], blockDepth)
                }
                overwrite(matchPos, ignoreLen)
              }
              if (blockDepth < blocksCnt && !continuation) {
                newContext = !0;
                do {
                  var startTag = blocks.pop();
                  addEndTag(startTag.getName(), textBoundary, 0).pairWith(startTag)
                }
                while (blockDepth < --blocksCnt);
              }
              if (blockDepth > blocksCnt && !lineIsEmpty) {
                newContext = !0;
                do {
                  var tagName = (blockMarks[blocksCnt] === '>!') ? 'SPOILER' : 'QUOTE';
                  blocks.push(addStartTag(tagName, matchPos, 0, -999))
                }
                while (blockDepth > ++blocksCnt);
              }
              var indentWidth = 0,
                indentPos = 0;
              if (m[2] && !codeFence) {
                indentStr = m[2];
                indentLen = indentStr.length;
                do {
                  if (indentStr[indentPos] === ' ') {
                    ++indentWidth
                  } else {
                    indentWidth = (indentWidth + 4) & ~3
                  }
                }
                while (++indentPos < indentLen && indentWidth < codeIndent);
              }
              if (codeTag && !codeFence && indentWidth < codeIndent && !lineIsEmpty) {
                newContext = !0
              }
              if (newContext) {
                newContext = !1;
                if (codeTag) {
                  if (textBoundary > codeTag.getPos()) {
                    overwrite(codeTag.getPos(), textBoundary - codeTag.getPos());
                    codeTag.pairWith(addEndTag('CODE', textBoundary, 0, -1))
                  } else {
                    codeTag.invalidate()
                  }
                  codeTag = null;
                  codeFence = null
                }
                lists.forEach(function (list) {
                  closeList(list, textBoundary)
                });
                lists = [];
                listsCnt = 0;
                if (matchPos) {
                  markBoundary(matchPos - 1)
                }
              }
              if (indentWidth >= codeIndent) {
                if (codeTag || !continuation) {
                  ignoreLen = (m[1] || '').length + indentPos;
                  if (!codeTag) {
                    codeTag = addStartTag('CODE', matchPos + ignoreLen, 0, -999)
                  }
                  m = {}
                }
              } else {
                var hasListItem = !!m[4];
                if (!indentWidth && !continuation && !hasListItem) {
                  listIndex = -1
                } else if (continuation && !hasListItem) {
                  listIndex = listsCnt - 1
                } else if (!listsCnt) {
                  listIndex = (hasListItem) ? 0 : -1
                } else {
                  listIndex = 0;
                  while (listIndex < listsCnt && indentWidth > lists[listIndex].maxIndent) {
                    ++listIndex
                  }
                }
                while (listIndex < listsCnt - 1) {
                  closeList(lists.pop(), textBoundary);
                  --listsCnt
                }
                if (listIndex === listsCnt && !hasListItem) {
                  --listIndex
                }
                if (hasListItem && listIndex >= 0) {
                  breakParagraph = !0;
                  tagPos = matchPos + ignoreLen + indentPos;
                  tagLen = m[4].length;
                  var itemTag = addStartTag('LI', tagPos, tagLen);
                  overwrite(tagPos, tagLen);
                  if (listIndex < listsCnt) {
                    addEndTag('LI', textBoundary, 0).pairWith(lists[listIndex].itemTag);
                    lists[listIndex].itemTag = itemTag;
                    lists[listIndex].itemTags.push(itemTag)
                  } else {
                    ++listsCnt;
                    if (listIndex) {
                      minIndent = lists[listIndex - 1].maxIndent + 1;
                      maxIndent = Math.max(minIndent, listIndex * 4)
                    } else {
                      minIndent = 0;
                      maxIndent = indentWidth
                    }
                    var listTag = addStartTag('LIST', tagPos, 0);
                    if (m[4].indexOf('.') > -1) {
                      listTag.setAttribute('type', 'decimal');
                      var start = +m[4];
                      if (start !== 1) {
                        listTag.setAttribute('start', start)
                      }
                    }
                    lists.push({
                      listTag: listTag,
                      itemTag: itemTag,
                      itemTags: [itemTag],
                      minIndent: minIndent,
                      maxIndent: maxIndent,
                      tight: !0
                    })
                  }
                }
                if (listsCnt && !continuation && !lineIsEmpty) {
                  if (lists[0].itemTags.length > 1 || !hasListItem) {
                    lists.forEach(function (list) {
                      list.tight = !1
                    })
                  }
                }
                codeIndent = (listsCnt + 1) * 4
              }
              if (m[5]) {
                if (m[5][0] === '#') {
                  startLen = m[5].length;
                  startPos = matchPos + matchLen - startLen;
                  endLen = getAtxHeaderEndTagLen(matchPos + matchLen, lfPos);
                  endPos = lfPos - endLen;
                  addTagPair('H' + /#{1,6}/.exec(m[5])[0].length, startPos, startLen, endPos, endLen);
                  markBoundary(startPos);
                  markBoundary(lfPos);
                  if (continuation) {
                    breakParagraph = !0
                  }
                } else if (m[5][0] === '`' || m[5][0] === '~') {
                  tagPos = matchPos + ignoreLen;
                  tagLen = lfPos - tagPos;
                  if (codeTag && m[5] === codeFence) {
                    codeTag.pairWith(addEndTag('CODE', tagPos, tagLen, -1));
                    addIgnoreTag(textBoundary, tagPos - textBoundary);
                    overwrite(codeTag.getPos(), tagPos + tagLen - codeTag.getPos());
                    codeTag = null;
                    codeFence = null
                  } else if (!codeTag) {
                    codeTag = addStartTag('CODE', tagPos, tagLen);
                    codeFence = m[5].replace(/[^`~]+/, '');
                    codeTag.setAttribute('blockDepth', blockDepth);
                    addIgnoreTag(tagPos + tagLen, 1);
                    var lang = m[5].replace(/^[`~\s]*/, '').replace(/\s+$/, '');
                    if (lang !== '') {
                      codeTag.setAttribute('lang', lang)
                    }
                  }
                }
              } else if (m[3] && !listsCnt && text[matchPos + matchLen] !== "\x17") {
                addSelfClosingTag('HR', matchPos + ignoreLen, matchLen - ignoreLen);
                breakParagraph = !0;
                markBoundary(lfPos)
              } else if (setextLines[lfPos] && setextLines[lfPos].blockDepth === blockDepth && !lineIsEmpty && !listsCnt && !codeTag) {
                addTagPair(setextLines[lfPos].tagName, matchPos + ignoreLen, 0, setextLines[lfPos].endPos, setextLines[lfPos].endLen);
                markBoundary(setextLines[lfPos].endPos + setextLines[lfPos].endLen)
              }
              if (breakParagraph) {
                addParagraphBreak(textBoundary);
                markBoundary(textBoundary)
              }
              if (!lineIsEmpty) {
                textBoundary = lfPos
              }
              if (ignoreLen) {
                addIgnoreTag(matchPos, ignoreLen, 1000)
              }
            })
          }

          function closeList(list, textBoundary) {
            addEndTag('LIST', textBoundary, 0).pairWith(list.listTag);
            addEndTag('LI', textBoundary, 0).pairWith(list.itemTag);
            if (list.tight) {
              list.itemTags.forEach(function (itemTag) {
                itemTag.removeFlags(RULE_CREATE_PARAGRAPHS)
              })
            }
          }

          function computeBlockIgnoreLen(str, maxBlockDepth) {
            var remaining = str;
            while (--maxBlockDepth >= 0) {
              remaining = remaining.replace(/^ *>!? ?/, '')
            }
            return str.length - remaining.length
          }

          function getAtxHeaderEndTagLen(startPos, endPos) {
            var content = text.substr(startPos, endPos - startPos),
              m = /[ \t]*#*[ \t]*$/.exec(content);
            return m[0].length
          }

          function getBlockMarks(str) {
            var blockMarks = [],
              regexp = />!?/g,
              m;
            while (m = regexp.exec(str)) {
              blockMarks.push(m[0])
            }
            return blockMarks
          }

          function matchSetextLines() {
            if (text.indexOf('-') === -1 && text.indexOf('=') === -1) {
              return
            }
            var m, regexp = /^(?=[-=>])(?:>!? ?)*(?=[-=])(?:-+|=+) *$/gm;
            while (m = regexp.exec(text)) {
              var match = m[0],
                matchPos = m.index;
              var endPos = matchPos - 1;
              while (endPos > 0 && text[endPos - 1] === ' ') {
                --endPos
              }
              setextLines[matchPos - 1] = {
                endLen: matchPos + match.length - endPos,
                endPos: endPos,
                blockDepth: match.length - match.replace(/>/g, '').length,
                tagName: (match[0] === '=') ? 'H1' : 'H2'
              }
            }
          }
          parse()
        })();
        (function () {
          function parse() {
            if (text.indexOf(']:') < 0) {
              return
            }
            var m, regexp = /^\x1A* {0,3}\[([^\x17\]]+)\]: *([^[\s\x17]+ *(?:"[^\x17]*?"|'[^\x17]*?'|\([^\x17)]*\))?) *(?=$|\x17)\n?/gm;
            while (m = regexp.exec(text)) {
              addIgnoreTag(m.index, m[0].length);
              var id = m[1].toLowerCase();
              if (!linkReferences[id]) {
                hasReferences = !0;
                linkReferences[id] = m[2]
              }
            }
          }
          parse()
        })();
        (function () {
          function parse() {
            var markers = getInlineCodeMarkers(),
              i = -1,
              cnt = markers.length;
            while (++i < (cnt - 1)) {
              var pos = markers[i].next,
                j = i;
              if (text[markers[i].pos] !== '`') {
                ++markers[i].pos;
                --markers[i].len
              }
              while (++j < cnt && markers[j].pos === pos) {
                if (markers[j].len === markers[i].len) {
                  addInlineCodeTags(markers[i], markers[j]);
                  i = j;
                  break
                }
                pos = markers[j].next
              }
            }
          }

          function addInlineCodeTags(left, right) {
            var startPos = left.pos,
              startLen = left.len + left.trimAfter,
              endPos = right.pos - right.trimBefore,
              endLen = right.len + right.trimBefore;
            addTagPair('C', startPos, startLen, endPos, endLen);
            overwrite(startPos, endPos + endLen - startPos)
          }

          function getInlineCodeMarkers() {
            var pos = text.indexOf('`');
            if (pos < 0) {
              return []
            }
            var regexp = /(`+)(\s*)[^\x17`]*/g,
              trimNext = 0,
              markers = [],
              _text = text.replace(/\x1BD/g, '\\`'),
              m;
            regexp.lastIndex = pos;
            while (m = regexp.exec(_text)) {
              markers.push({
                pos: m.index,
                len: m[1].length,
                trimBefore: trimNext,
                trimAfter: m[2].length,
                next: m.index + m[0].length
              });
              trimNext = m[0].length - m[0].replace(/\s+$/, '').length
            }
            return markers
          }
          parse()
        })();
        (function () {
          function parse() {
            var pos = text.indexOf('![');
            if (pos === -1) {
              return
            }
            if (text.indexOf('](', pos) > 0) {
              parseInlineImages()
            }
            if (hasReferences) {
              parseReferenceImages()
            }
          }

          function addImageTag(startPos, endPos, endLen, linkInfo, alt) {
            var tag = addTagPair('IMG', startPos, 2, endPos, endLen);
            setLinkAttributes(tag, linkInfo, 'src');
            tag.setAttribute('alt', decode(alt));
            overwrite(startPos, endPos + endLen - startPos)
          }

          function parseInlineImages() {
            var m, regexp = /!\[(?:[^\x17[\]]|\[[^\x17[\]]*\])*\]\(( *(?:[^\x17\s()]|\([^\x17\s()]*\))*(?=[ )]) *(?:"[^\x17]*?"|'[^\x17]*?'|\([^\x17)]*\))? *)\)/g;
            while (m = regexp.exec(text)) {
              var linkInfo = m[1],
                startPos = m.index,
                endLen = 3 + linkInfo.length,
                endPos = startPos + m[0].length - endLen,
                alt = m[0].substr(2, m[0].length - endLen - 2);
              addImageTag(startPos, endPos, endLen, linkInfo, alt)
            }
          }

          function parseReferenceImages() {
            var m, regexp = /!\[((?:[^\x17[\]]|\[[^\x17[\]]*\])*)\](?: ?\[([^\x17[\]]+)\])?/g;
            while (m = regexp.exec(text)) {
              var startPos = m.index,
                endPos = startPos + 2 + m[1].length,
                endLen = 1,
                alt = m[1],
                id = alt;
              if (m[2] > '' && linkReferences[m[2]]) {
                endLen = m[0].length - alt.length - 2;
                id = m[2]
              } else if (!linkReferences[id]) {
                continue
              }
              addImageTag(startPos, endPos, endLen, linkReferences[id], alt)
            }
          }
          parse()
        })();
        (function () {
          function parse() {
            parseInlineMarkup('>!', />![^\x17]+?!</g, 'ISPOILER');
            parseInlineMarkup('||', /\|\|[^\x17]+?\|\|/g, 'ISPOILER')
          }
          parse()
        })();
        (function () {
          function parse() {
            if (text.indexOf('](') !== -1) {
              parseInlineLinks()
            }
            if (text.indexOf('<') !== -1) {
              parseAutomaticLinks()
            }
            if (hasReferences) {
              parseReferenceLinks()
            }
          }

          function addLinkTag(startPos, endPos, endLen, linkInfo) {
            var priority = (endLen === 1) ? 1 : -1;
            var tag = addTagPair('URL', startPos, 1, endPos, endLen, priority);
            setLinkAttributes(tag, linkInfo, 'url');
            overwrite(startPos, 1);
            overwrite(endPos, endLen)
          }

          function getLabels() {
            var labels = {},
              m, regexp = /\[((?:[^\x17[\]]|\[[^\x17[\]]*\])*)\]/g;
            while (m = regexp.exec(text)) {
              labels[m.index] = m[1].toLowerCase()
            }
            return labels
          }

          function parseAutomaticLinks() {
            var m, regexp = /<[-+.\w]+([:@])[^\x17\s>]+?(?:>|\x1B7)/g;
            while (m = regexp.exec(text)) {
              var content = decode(m[0].replace(/\x1B/g, "\\\x1B")).replace(/^<(.+)>$/, '$1'),
                startPos = m.index,
                endPos = startPos + m[0].length - 1,
                tagName = (m[1] === ':') ? 'URL' : 'EMAIL',
                attrName = tagName.toLowerCase();
              addTagPair(tagName, startPos, 1, endPos, 1).setAttribute(attrName, content)
            }
          }

          function parseInlineLinks() {
            var m, regexp = /\[(?:[^\x17[\]]|\[[^\x17[\]]*\])*\]\(( *(?:[^\x17\s()]|\([^\x17\s()]*\))*(?=[ )]) *(?:"[^\x17]*?"|'[^\x17]*?'|\([^\x17)]*\))? *)\)/g;
            while (m = regexp.exec(text)) {
              var linkInfo = m[1],
                startPos = m.index,
                endLen = 3 + linkInfo.length,
                endPos = startPos + m[0].length - endLen;
              addLinkTag(startPos, endPos, endLen, linkInfo)
            }
          }

          function parseReferenceLinks() {
            var labels = getLabels(),
              startPos;
            for (startPos in labels) {
              var id = labels[startPos],
                labelPos = +startPos + 2 + id.length,
                endPos = labelPos - 1,
                endLen = 1;
              if (text[labelPos] === ' ') {
                ++labelPos
              }
              if (labels[labelPos] > '' && linkReferences[labels[labelPos]]) {
                id = labels[labelPos];
                endLen = labelPos + 2 + id.length - endPos
              }
              if (linkReferences[id]) {
                addLinkTag(+startPos, endPos, endLen, linkReferences[id])
              }
            }
          }
          parse()
        })();
        (function () {
          function parse() {
            parseInlineMarkup('~~', /~~[^\x17]+?~~(?!~)/g, 'DEL')
          }
          parse()
        })();
        (function () {
          function parse() {
            parseAbstractScript('SUB', '~', /~[^\x17\s!"#$%&\'()*+,\-.\/:;<=>?@[\]^_`{}|~]+~?/g, /~\([^\x17()]+\)/g)
          }
          parse()
        })();
        (function () {
          function parse() {
            parseAbstractScript('SUP', '^', /\^[^\x17\s!"#$%&\'()*+,\-.\/:;<=>?@[\]^_`{}|~]+\^?/g, /\^\([^\x17()]+\)/g)
          }
          parse()
        })();
        (function () {
          var closeEm;
          var closeStrong;
          var emPos;
          var emEndPos;
          var remaining;
          var strongPos;
          var strongEndPos;

          function parse() {
            parseEmphasisByCharacter('*', /\*+/g);
            parseEmphasisByCharacter('_', /_+/g)
          }

          function adjustEndingPositions() {
            if (closeEm && closeStrong) {
              if (emPos < strongPos) {
                emEndPos += 2
              } else {
                ++strongEndPos
              }
            }
          }

          function adjustStartingPositions() {
            if (emPos >= 0 && emPos === strongPos) {
              if (closeEm) {
                emPos += 2
              } else {
                ++strongPos
              }
            }
          }

          function closeSpans() {
            if (closeEm) {
              --remaining;
              addTagPair('EM', emPos, 1, emEndPos, 1);
              emPos = -1
            }
            if (closeStrong) {
              remaining -= 2;
              addTagPair('STRONG', strongPos, 2, strongEndPos, 2);
              strongPos = -1
            }
          }

          function getEmphasisByBlock(regexp, pos) {
            var block = [],
              blocks = [],
              breakPos = text.indexOf("\x17", pos),
              m;
            regexp.lastIndex = pos;
            while (m = regexp.exec(text)) {
              var matchPos = m.index,
                matchLen = m[0].length;
              if (matchPos > breakPos) {
                blocks.push(block);
                block = [];
                breakPos = text.indexOf("\x17", matchPos)
              }
              if (!ignoreEmphasis(matchPos, matchLen)) {
                block.push([matchPos, matchLen])
              }
            }
            blocks.push(block);
            return blocks
          }

          function ignoreEmphasis(pos, len) {
            return (text.charAt(pos) === '_' && len === 1 && isSurroundedByAlnum(pos, len))
          }

          function openSpans(pos) {
            if (remaining & 1) {
              emPos = pos - remaining
            }
            if (remaining & 2) {
              strongPos = pos - remaining
            }
          }

          function parseEmphasisByCharacter(character, regexp) {
            var pos = text.indexOf(character);
            if (pos === -1) {
              return
            }
            getEmphasisByBlock(regexp, pos).forEach(processEmphasisBlock)
          }

          function processEmphasisBlock(block) {
            emPos = -1, strongPos = -1;
            block.forEach(function (pair) {
              processEmphasisMatch(pair[0], pair[1])
            })
          }

          function processEmphasisMatch(matchPos, matchLen) {
            var canOpen = !isBeforeWhitespace(matchPos + matchLen - 1),
              canClose = !isAfterWhitespace(matchPos),
              closeLen = (canClose) ? Math.min(matchLen, 3) : 0;
            closeEm = !!(closeLen & 1) && emPos >= 0;
            closeStrong = !!(closeLen & 2) && strongPos >= 0;
            emEndPos = matchPos;
            strongEndPos = matchPos;
            remaining = matchLen;
            adjustStartingPositions();
            adjustEndingPositions();
            closeSpans();
            remaining = (canOpen) ? Math.min(remaining, 3) : 0;
            openSpans(matchPos + matchLen)
          }
          parse()
        })();
        (function () {
          function parse() {
            var pos = text.indexOf("  \n");
            while (pos > 0) {
              addBrTag(pos + 2).cascadeInvalidationTo(addVerbatim(pos + 2, 1));
              pos = text.indexOf("  \n", pos + 3)
            }
          }
          parse()
        })()
      }
    },
    "Preg": {
      parser: function (text, matches) {
        var config = {
          generics: [
            ["USERMENTION", /\B@([a-z0-9_-]+)(?!#)/ig, 0, ["", "username"]],
            ["POSTMENTION", /\B@([a-z0-9_-]+)#(\d+)/ig, 0, ["", "username", "id"]]
          ]
        };
        config.generics.forEach(function (entry) {
          var tagName = entry[0],
            regexp = entry[1],
            passthroughIdx = entry[2],
            map = entry[3],
            m;
          regexp.lastIndex = 0;
          while (m = regexp.exec(text)) {
            var startTagPos = m.index,
              matchLen = m[0].length,
              tag;
            if (HINT.PREG_HAS_PASSTHROUGH && passthroughIdx && m[passthroughIdx] !== '') {
              var contentPos = text.indexOf(m[passthroughIdx], startTagPos),
                contentLen = m[passthroughIdx].length,
                startTagLen = contentPos - startTagPos,
                endTagPos = contentPos + contentLen,
                endTagLen = matchLen - (startTagLen + contentLen);
              tag = addTagPair(tagName, startTagPos, startTagLen, endTagPos, endTagLen, -100)
            } else {
              tag = addSelfClosingTag(tagName, startTagPos, matchLen, -100)
            }
            map.forEach(function (attrName, i) {
              if (attrName && typeof m[i] !== 'undefined') {
                tag.setAttribute(attrName, m[i])
              }
            })
          }
        })
      }
    }
  };
  var pos;
  var registeredVars = {
    "urlConfig": {
      allowedSchemes: /^https?$/i
    }
  };
  var rootContext = {
    allowed: oF6AF222C,
    flags: 40
  };
  var tagsConfig = {
    "B": o5B6ED7AA,
    "C": {
      allowed: o553403F9,
      attributes: {},
      bitNumber: 2,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: {
        flags: 66
      },
      tagLimit: 5000
    },
    "CENTER": {
      allowed: oF6AF222C,
      attributes: {},
      bitNumber: 3,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: o1AF69066,
      tagLimit: 5000
    },
    "CODE": {
      allowed: o553403F9,
      attributes: {
        "lang": {
          filterChain: [function (attrValue, attrName) {
            return RegexpFilter.filter(attrValue, /^[-0-9A-Za-z_]+$/)
          }],
          required: !1
        }
      },
      bitNumber: 3,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: {
        closeParent: oD4869BFF,
        flags: 4436,
        fosterParent: oD4869BFF
      },
      tagLimit: 5000
    },
    "COLOR": {
      allowed: oD363F9C5,
      attributes: {
        "color": {
          filterChain: [function (attrValue, attrName) {
            return RegexpFilter.filter(attrValue, /^(?:#(?:(?:[0-9a-f]{3}){1,2}|(?:[0-9a-f]{4}){1,2})|rgb\(\d{1,3}, *\d{1,3}, *\d{1,3}\)|rgba\(\d{1,3}, *\d{1,3}, *\d{1,3}, *\d*(?:\.\d+)?\)|[a-z]+)$/i)
          }],
          required: !0
        }
      },
      bitNumber: 2,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: oB565876D,
      tagLimit: 5000
    },
    "DEL": {
      allowed: oF6AF222C,
      attributes: {},
      bitNumber: 2,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: {
        flags: 512
      },
      tagLimit: 5000
    },
    "E": {
      allowed: oC06C5BF5,
      attributes: {},
      bitNumber: 5,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: oA80287CC,
      tagLimit: 5000
    },
    "EM": o5B6ED7AA,
    "EMAIL": {
      allowed: o3C789569,
      attributes: {
        "email": {
          filterChain: [function (attrValue, attrName) {
            return EmailFilter.filter(attrValue)
          }],
          required: !0
        }
      },
      bitNumber: 1,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: o939F1698,
      tagLimit: 5000
    },
    "ESC": {
      allowed: o553403F9,
      attributes: {},
      bitNumber: 0,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: {
        flags: 1616
      },
      tagLimit: 5000
    },
    "H1": oE13673F5,
    "H2": oE13673F5,
    "H3": oE13673F5,
    "H4": oE13673F5,
    "H5": oE13673F5,
    "H6": oE13673F5,
    "HR": {
      allowed: oC06C5BF5,
      attributes: {},
      bitNumber: 3,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: {
        closeParent: oD4869BFF,
        flags: 3349
      },
      tagLimit: 5000
    },
    "I": o5B6ED7AA,
    "IMG": {
      allowed: oC06C5BF5,
      attributes: {
        "alt": o1BC3EAF4,
        "height": o75AB9259,
        "src": oF307D35C,
        "title": o1BC3EAF4,
        "width": o75AB9259
      },
      bitNumber: 2,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: oA80287CC,
      tagLimit: 5000
    },
    "ISPOILER": o8F52476D,
    "LI": {
      allowed: oF6AF222C,
      attributes: {},
      bitNumber: 4,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: {
        closeParent: {
          "B": 1,
          "C": 1,
          "COLOR": 1,
          "EM": 1,
          "EMAIL": 1,
          "I": 1,
          "LI": 1,
          "S": 1,
          "SIZE": 1,
          "STRONG": 1,
          "U": 1,
          "URL": 1
        },
        flags: 264,
        fosterParent: oD4869BFF
      },
      tagLimit: 5000
    },
    "LIST": {
      allowed: [32529],
      attributes: {
        "start": o75AB9259,
        "type": {
          filterChain: [function (attrValue, attrName) {
            return HashmapFilter.filter(attrValue, {
              "A": "upper-alpha",
              "I": "upper-roman",
              "a": "lower-alpha",
              "i": "lower-roman",
              "1": "decimal"
            }, !1)
          }, function (attrValue, attrName) {
            return RegexpFilter.filter(attrValue, /^[- +,.0-9A-Za-z_]+$/)
          }],
          required: !1
        }
      },
      bitNumber: 3,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: {
        closeParent: oD4869BFF,
        createChild: ["LI"],
        flags: 3460,
        fosterParent: oD4869BFF
      },
      tagLimit: 5000
    },
    "POSTMENTION": {
      allowed: oC06C5BF5,
      attributes: {
        "discussionid": o6CB0A318,
        "displayname": o02D8DBB5,
        "id": o6CB0A318,
        "number": o6CB0A318,
        "username": o02D8DBB5
      },
      bitNumber: 1,
      filterChain: [function (tag, tagConfig) {
        return (function (tag) {
          return flarum.extensions["flarum-mentions"].filterPostMentions(tag)
        })(tag)
      }, o118B31AC],
      nestingLimit: 10,
      rules: oA80287CC,
      tagLimit: 5000
    },
    "QUOTE": {
      allowed: oF6AF222C,
      attributes: {
        "author": o1BC3EAF4
      },
      bitNumber: 3,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: o1AF69066,
      tagLimit: 5000
    },
    "S": o5B6ED7AA,
    "SIZE": {
      allowed: oD363F9C5,
      attributes: {
        "size": {
          filterChain: [function (attrValue, attrName) {
            return NumericFilter.filterRange(attrValue, 8, 36, logger)
          }],
          required: !0
        }
      },
      bitNumber: 2,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: oB565876D,
      tagLimit: 5000
    },
    "SPOILER": {
      allowed: oF6AF222C,
      attributes: {},
      bitNumber: 6,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: o1AF69066,
      tagLimit: 5000
    },
    "STRONG": o5B6ED7AA,
    "SUB": o8F52476D,
    "SUP": o8F52476D,
    "U": o5B6ED7AA,
    "URL": {
      allowed: o3C789569,
      attributes: {
        "title": o1BC3EAF4,
        "url": oF307D35C
      },
      bitNumber: 1,
      filterChain: o1B4F5B52,
      nestingLimit: 10,
      rules: o939F1698,
      tagLimit: 5000
    },
    "USERMENTION": {
      allowed: oC06C5BF5,
      attributes: {
        "displayname": o02D8DBB5,
        "id": o6CB0A318,
        "username": o02D8DBB5
      },
      bitNumber: 1,
      filterChain: [function (tag, tagConfig) {
        return (function (tag) {
          return flarum.extensions["flarum-mentions"].filterUserMentions(tag)
        })(tag)
      }, o118B31AC],
      nestingLimit: 10,
      rules: oA80287CC,
      tagLimit: 5000
    }
  };
  var tagStack;
  var tagStackIsSorted;
  var text;
  var textLen;
  var uid = 0;
  var wsPos;

  function disableTag(tagName) {
    if (tagsConfig[tagName]) {
      copyTagConfig(tagName).isDisabled = !0
    }
  }

  function enableTag(tagName) {
    if (tagsConfig[tagName]) {
      copyTagConfig(tagName).isDisabled = !1
    }
  }

  function getLogger() {
    return logger
  }

  function parse(_text) {
    reset(_text);
    var _uid = uid;
    executePluginParsers();
    processTags();
    finalizeOutput();
    if (uid !== _uid) {
      throw 'The parser has been reset during execution'
    }
    if (currentFixingCost > maxFixingCost) {
      logger.warn('Fixing cost limit exceeded')
    }
    return output
  }

  function reset(_text) {
    _text = _text.replace(/\r\n?/g, "\n");
    _text = _text.replace(/[\x00-\x08\x0B\x0C\x0E-\x1F]+/g, '');
    logger.clear();
    cntOpen = {};
    cntTotal = {};
    currentFixingCost = 0;
    currentTag = null;
    isRich = !1;
    namespaces = {};
    openTags = [];
    output = '';
    pos = 0;
    tagStack = [];
    tagStackIsSorted = !1;
    text = _text;
    textLen = text.length;
    wsPos = 0;
    context = rootContext;
    context.inParagraph = !1;
    ++uid
  }

  function setTagLimit(tagName, tagLimit) {
    if (tagsConfig[tagName]) {
      copyTagConfig(tagName).tagLimit = tagLimit
    }
  }

  function setNestingLimit(tagName, nestingLimit) {
    if (tagsConfig[tagName]) {
      copyTagConfig(tagName).nestingLimit = nestingLimit
    }
  }

  function copyTagConfig(tagName) {
    var tagConfig = {},
      k;
    for (k in tagsConfig[tagName]) {
      tagConfig[k] = tagsConfig[tagName][k]
    }
    return tagsConfig[tagName] = tagConfig
  }

  function encodeUnicodeSupplementaryCharacters() {
    output = output.replace(/[\uD800-\uDBFF][\uDC00-\uDFFF]/g, encodeUnicodeSupplementaryCharactersCallback)
  }

  function encodeUnicodeSupplementaryCharactersCallback(pair) {
    var cp = (pair.charCodeAt(0) << 10) + pair.charCodeAt(1) - 56613888;
    return '&#' + cp + ';'
  }

  function finalizeOutput() {
    var tmp;
    outputText(textLen, 0, !0);
    do {
      tmp = output;
      output = output.replace(/<([^ />]+)[^>]*><\/\1>/g, '')
    }
    while (output !== tmp);
    output = output.replace(/<\/i><i>/g, '');
    output = output.replace(/[\x00-\x08\x0B-\x1F]/g, '');
    encodeUnicodeSupplementaryCharacters();
    var tagName = (isRich) ? 'r' : 't';
    tmp = '<' + tagName;
    if (HINT.namespaces) {
      for (var prefix in namespaces) {
        tmp += ' xmlns:' + prefix + '="urn:s9e:TextFormatter:' + prefix + '"'
      }
    }
    output = tmp + '>' + output + '</' + tagName + '>'
  }

  function outputTag(tag) {
    isRich = !0;
    var tagName = tag.getName(),
      tagPos = tag.getPos(),
      tagLen = tag.getLen(),
      tagFlags = tag.getFlags(),
      skipBefore = 0,
      skipAfter = 0;
    if (HINT.RULE_IGNORE_WHITESPACE && (tagFlags & RULE_IGNORE_WHITESPACE)) {
      skipBefore = 1;
      skipAfter = (tag.isEndTag()) ? 2 : 1
    }
    var closeParagraph = !1;
    if (tag.isStartTag()) {
      if (HINT.RULE_BREAK_PARAGRAPH && (tagFlags & RULE_BREAK_PARAGRAPH)) {
        closeParagraph = !0
      }
    } else {
      closeParagraph = !0
    }
    outputText(tagPos, skipBefore, closeParagraph);
    var tagText = (tagLen) ? htmlspecialchars_noquotes(text.substr(tagPos, tagLen)) : '';
    if (tag.isStartTag()) {
      if (!HINT.RULE_BREAK_PARAGRAPH || !(tagFlags & RULE_BREAK_PARAGRAPH)) {
        outputParagraphStart(tagPos)
      }
      if (HINT.namespaces) {
        var colonPos = tagName.indexOf(':');
        if (colonPos > 0) {
          namespaces[tagName.substr(0, colonPos)] = 0
        }
      }
      output += '<' + tagName;
      var attributes = tag.getAttributes(),
        attributeNames = [];
      for (var attrName in attributes) {
        attributeNames.push(attrName)
      }
      attributeNames.sort(function (a, b) {
        return (a > b) ? 1 : -1
      });
      attributeNames.forEach(function (attrName) {
        output += ' ' + attrName + '="' + htmlspecialchars_compat(attributes[attrName].toString()).replace(/\n/g, '&#10;') + '"'
      });
      if (tag.isSelfClosingTag()) {
        if (tagLen) {
          output += '>' + tagText + '</' + tagName + '>'
        } else {
          output += '/>'
        }
      } else if (tagLen) {
        output += '><s>' + tagText + '</s>'
      } else {
        output += '>'
      }
    } else {
      if (tagLen) {
        output += '<e>' + tagText + '</e>'
      }
      output += '</' + tagName + '>'
    }
    pos = tagPos + tagLen;
    wsPos = pos;
    while (skipAfter && wsPos < textLen && text[wsPos] === "\n") {
      --skipAfter;
      ++wsPos
    }
  }

  function outputText(catchupPos, maxLines, closeParagraph) {
    if (closeParagraph) {
      if (!(context.flags & RULE_CREATE_PARAGRAPHS)) {
        closeParagraph = !1
      } else {
        maxLines = -1
      }
    }
    if (pos >= catchupPos) {
      if (closeParagraph) {
        outputParagraphEnd()
      }
    }
    if (wsPos > pos) {
      var skipPos = Math.min(catchupPos, wsPos);
      output += text.substr(pos, skipPos - pos);
      pos = skipPos;
      if (pos >= catchupPos) {
        if (closeParagraph) {
          outputParagraphEnd()
        }
      }
    }
    var catchupLen, catchupText;
    if (HINT.RULE_IGNORE_TEXT && context.flags & RULE_IGNORE_TEXT) {
      catchupLen = catchupPos - pos, catchupText = text.substr(pos, catchupLen);
      if (!/^[ \n\t]*$/.test(catchupText)) {
        catchupText = '<i>' + htmlspecialchars_noquotes(catchupText) + '</i>'
      }
      output += catchupText;
      pos = catchupPos;
      if (closeParagraph) {
        outputParagraphEnd()
      }
      return
    }
    var ignorePos = catchupPos,
      ignoreLen = 0;
    while (maxLines && --ignorePos >= pos) {
      var c = text[ignorePos];
      if (c !== ' ' && c !== "\n" && c !== "\t") {
        break
      }
      if (c === "\n") {
        --maxLines
      }
      ++ignoreLen
    }
    catchupPos -= ignoreLen;
    if (HINT.RULE_CREATE_PARAGRAPHS && context.flags & RULE_CREATE_PARAGRAPHS) {
      if (!context.inParagraph) {
        outputWhitespace(catchupPos);
        if (catchupPos > pos) {
          outputParagraphStart(catchupPos)
        }
      }
      var pbPos = text.indexOf("\n\n", pos);
      while (pbPos > -1 && pbPos < catchupPos) {
        outputText(pbPos, 0, !0);
        outputParagraphStart(catchupPos);
        pbPos = text.indexOf("\n\n", pos)
      }
    }
    if (catchupPos > pos) {
      catchupText = htmlspecialchars_noquotes(text.substr(pos, catchupPos - pos));
      if (HINT.RULE_ENABLE_AUTO_BR && (context.flags & RULES_AUTO_LINEBREAKS) === RULE_ENABLE_AUTO_BR) {
        catchupText = catchupText.replace(/\n/g, "<br/>\n")
      }
      output += catchupText
    }
    if (closeParagraph) {
      outputParagraphEnd()
    }
    if (ignoreLen) {
      output += text.substr(catchupPos, ignoreLen)
    }
    pos = catchupPos + ignoreLen
  }

  function outputBrTag(tag) {
    outputText(tag.getPos(), 0, !1);
    output += '<br/>'
  }

  function outputIgnoreTag(tag) {
    var tagPos = tag.getPos(),
      tagLen = tag.getLen();
    var ignoreText = text.substr(tagPos, tagLen);
    outputText(tagPos, 0, !1);
    output += '<i>' + htmlspecialchars_noquotes(ignoreText) + '</i>';
    isRich = !0;
    pos = tagPos + tagLen
  }

  function outputParagraphStart(maxPos) {
    if (!HINT.RULE_CREATE_PARAGRAPHS) {
      return
    }
    if (context.inParagraph || !(context.flags & RULE_CREATE_PARAGRAPHS)) {
      return
    }
    outputWhitespace(maxPos);
    if (pos < textLen) {
      output += '<p>';
      context.inParagraph = !0
    }
  }

  function outputParagraphEnd() {
    if (!context.inParagraph) {
      return
    }
    output += '</p>';
    context.inParagraph = !1
  }

  function outputVerbatim(tag) {
    var flags = context.flags;
    context.flags = tag.getFlags();
    outputText(currentTag.getPos() + currentTag.getLen(), 0, !1);
    context.flags = flags
  }

  function outputWhitespace(maxPos) {
    while (pos < maxPos && " \n\t".indexOf(text[pos]) > -1) {
      output += text[pos];
      ++pos
    }
  }

  function disablePlugin(pluginName) {
    if (plugins[pluginName]) {
      plugins[pluginName].isDisabled = !0
    }
  }

  function enablePlugin(pluginName) {
    if (plugins[pluginName]) {
      plugins[pluginName].isDisabled = !1
    }
  }

  function executePluginParser(pluginName) {
    var pluginConfig = plugins[pluginName];
    if (pluginConfig.quickMatch && text.indexOf(pluginConfig.quickMatch) < 0) {
      return
    }
    var matches = [];
    if (HINT.regexp && HINT.regexpLimit && typeof pluginConfig.regexp !== 'undefined' && typeof pluginConfig.regexpLimit !== 'undefined') {
      matches = getMatches(pluginConfig.regexp, pluginConfig.regexpLimit);
      if (!matches.length) {
        return
      }
    }
    getPluginParser(pluginName)(text, matches)
  }

  function executePluginParsers() {
    for (var pluginName in plugins) {
      if (!plugins[pluginName].isDisabled) {
        executePluginParser(pluginName)
      }
    }
  }

  function getMatches(regexp, limit) {
    regexp.lastIndex = 0;
    var matches = [],
      cnt = 0,
      m;
    while (++cnt <= limit && (m = regexp.exec(text))) {
      var pos = m.index,
        match = [
          [m[0], pos]
        ],
        i = 0;
      while (++i < m.length) {
        var str = m[i];
        if (str === undefined) {
          match.push(['', -1])
        } else {
          match.push([str, text.indexOf(str, pos)]);
          pos += str.length
        }
      }
      matches.push(match)
    }
    return matches
  }

  function getPluginParser(pluginName) {
    return plugins[pluginName].parser
  }

  function registerParser(pluginName, parser, regexp, limit) {
    if (!plugins[pluginName]) {
      plugins[pluginName] = {}
    }
    if (regexp) {
      plugins[pluginName].regexp = regexp;
      plugins[pluginName].limit = limit || Infinity
    }
    plugins[pluginName].parser = parser
  }

  function closeAncestor(tag) {
    if (!HINT.closeAncestor) {
      return !1
    }
    if (openTags.length) {
      var tagName = tag.getName(),
        tagConfig = tagsConfig[tagName];
      if (tagConfig.rules.closeAncestor) {
        var i = openTags.length;
        while (--i >= 0) {
          var ancestor = openTags[i],
            ancestorName = ancestor.getName();
          if (tagConfig.rules.closeAncestor[ancestorName]) {
            ++currentFixingCost;
            tagStack.push(tag);
            addMagicEndTag(ancestor, tag.getPos(), tag.getSortPriority() - 1);
            return !0
          }
        }
      }
    }
    return !1
  }

  function closeParent(tag) {
    if (!HINT.closeParent) {
      return !1
    }
    if (openTags.length) {
      var tagName = tag.getName(),
        tagConfig = tagsConfig[tagName];
      if (tagConfig.rules.closeParent) {
        var parent = openTags[openTags.length - 1],
          parentName = parent.getName();
        if (tagConfig.rules.closeParent[parentName]) {
          ++currentFixingCost;
          tagStack.push(tag);
          addMagicEndTag(parent, tag.getPos(), tag.getSortPriority() - 1);
          return !0
        }
      }
    }
    return !1
  }

  function createChild(tag) {
    if (!HINT.createChild) {
      return
    }
    var tagConfig = tagsConfig[tag.getName()];
    if (tagConfig.rules.createChild) {
      var priority = -1000,
        _text = text.substr(pos),
        tagPos = pos + _text.length - _text.replace(/^[ \n\r\t]+/, '').length;
      tagConfig.rules.createChild.forEach(function (tagName) {
        addStartTag(tagName, tagPos, 0, ++priority)
      })
    }
  }

  function fosterParent(tag) {
    if (!HINT.fosterParent) {
      return !1
    }
    if (openTags.length) {
      var tagName = tag.getName(),
        tagConfig = tagsConfig[tagName];
      if (tagConfig.rules.fosterParent) {
        var parent = openTags[openTags.length - 1],
          parentName = parent.getName();
        if (tagConfig.rules.fosterParent[parentName]) {
          if (parentName !== tagName && currentFixingCost < maxFixingCost) {
            addFosterTag(tag, parent)
          }
          tagStack.push(tag);
          addMagicEndTag(parent, tag.getPos(), tag.getSortPriority() - 1);
          currentFixingCost += 4;
          return !0
        }
      }
    }
    return !1
  }

  function requireAncestor(tag) {
    if (!HINT.requireAncestor) {
      return !1
    }
    var tagName = tag.getName(),
      tagConfig = tagsConfig[tagName];
    if (tagConfig.rules.requireAncestor) {
      var i = tagConfig.rules.requireAncestor.length;
      while (--i >= 0) {
        var ancestorName = tagConfig.rules.requireAncestor[i];
        if (cntOpen[ancestorName]) {
          return !1
        }
      }
      logger.err('Tag requires an ancestor', {
        'requireAncestor': tagConfig.rules.requireAncestor.join(', '),
        'tag': tag
      });
      return !0
    }
    return !1
  }

  function addFosterTag(tag, fosterTag) {
    var coords = getMagicStartCoords(tag.getPos() + tag.getLen()),
      childPos = coords[0],
      childPrio = coords[1];
    var childTag = addCopyTag(fosterTag, childPos, 0, childPrio);
    tag.cascadeInvalidationTo(childTag)
  }

  function addMagicEndTag(startTag, tagPos, prio) {
    var tagName = startTag.getName();
    if (HINT.RULE_IGNORE_WHITESPACE && ((currentTag.getFlags() | startTag.getFlags()) & RULE_IGNORE_WHITESPACE)) {
      tagPos = getMagicEndPos(tagPos)
    }
    var endTag = addEndTag(tagName, tagPos, 0, prio || 0);
    endTag.pairWith(startTag);
    return endTag
  }

  function getMagicEndPos(tagPos) {
    while (tagPos > pos && WHITESPACE.indexOf(text[tagPos - 1]) > -1) {
      --tagPos
    }
    return tagPos
  }

  function getMagicStartCoords(tagPos) {
    var nextPos, nextPrio, nextTag, prio;
    if (!tagStack.length) {
      nextPos = textLen + 1;
      nextPrio = 0
    } else {
      nextTag = tagStack[tagStack.length - 1];
      nextPos = nextTag.getPos();
      nextPrio = nextTag.getSortPriority()
    }
    while (tagPos < nextPos && WHITESPACE.indexOf(text[tagPos]) > -1) {
      ++tagPos
    }
    prio = (tagPos === nextPos) ? nextPrio - 1 : 0;
    return [tagPos, prio]
  }

  function isFollowedByClosingTag(tag) {
    return (!tagStack.length) ? !1 : tagStack[tagStack.length - 1].canClose(tag)
  }

  function processTags() {
    if (!tagStack.length) {
      return
    }
    for (var tagName in tagsConfig) {
      cntOpen[tagName] = 0;
      cntTotal[tagName] = 0
    }
    do {
      while (tagStack.length) {
        if (!tagStackIsSorted) {
          sortTags()
        }
        currentTag = tagStack.pop();
        processCurrentTag()
      }
      openTags.forEach(function (startTag) {
        addMagicEndTag(startTag, textLen)
      })
    }
    while (tagStack.length);
  }

  function processCurrentTag() {
    if ((context.flags & RULE_IGNORE_TAGS) && !currentTag.canClose(openTags[openTags.length - 1]) && !currentTag.isSystemTag()) {
      currentTag.invalidate()
    }
    var tagPos = currentTag.getPos(),
      tagLen = currentTag.getLen();
    if (pos > tagPos && !currentTag.isInvalid()) {
      var startTag = currentTag.getStartTag();
      if (startTag && openTags.indexOf(startTag) >= 0) {
        addEndTag(startTag.getName(), pos, Math.max(0, tagPos + tagLen - pos)).pairWith(startTag);
        return
      }
      if (currentTag.isIgnoreTag()) {
        var ignoreLen = tagPos + tagLen - pos;
        if (ignoreLen > 0) {
          addIgnoreTag(pos, ignoreLen);
          return
        }
      }
      currentTag.invalidate()
    }
    if (currentTag.isInvalid()) {
      return
    }
    if (currentTag.isIgnoreTag()) {
      outputIgnoreTag(currentTag)
    } else if (currentTag.isBrTag()) {
      if (!HINT.RULE_PREVENT_BR || !(context.flags & RULE_PREVENT_BR)) {
        outputBrTag(currentTag)
      }
    } else if (currentTag.isParagraphBreak()) {
      outputText(currentTag.getPos(), 0, !0)
    } else if (currentTag.isVerbatim()) {
      outputVerbatim(currentTag)
    } else if (currentTag.isStartTag()) {
      processStartTag(currentTag)
    } else {
      processEndTag(currentTag)
    }
  }

  function processStartTag(tag) {
    var tagName = tag.getName(),
      tagConfig = tagsConfig[tagName];
    if (cntTotal[tagName] >= tagConfig.tagLimit) {
      logger.err('Tag limit exceeded', {
        'tag': tag,
        'tagName': tagName,
        'tagLimit': tagConfig.tagLimit
      });
      tag.invalidate();
      return
    }
    filterTag(tag);
    if (tag.isInvalid()) {
      return
    }
    if (currentFixingCost < maxFixingCost) {
      if (fosterParent(tag) || closeParent(tag) || closeAncestor(tag)) {
        return
      }
    }
    if (cntOpen[tagName] >= tagConfig.nestingLimit) {
      logger.err('Nesting limit exceeded', {
        'tag': tag,
        'tagName': tagName,
        'nestingLimit': tagConfig.nestingLimit
      });
      tag.invalidate();
      return
    }
    if (!tagIsAllowed(tagName)) {
      var msg = 'Tag is not allowed in this context',
        context = {
          'tag': tag,
          'tagName': tagName
        };
      if (tag.getLen() > 0) {
        logger.warn(msg, context)
      } else {
        logger.debug(msg, context)
      }
      tag.invalidate();
      return
    }
    if (requireAncestor(tag)) {
      tag.invalidate();
      return
    }
    if (HINT.RULE_AUTO_CLOSE && tag.getFlags() & RULE_AUTO_CLOSE && !tag.isSelfClosingTag() && !tag.getEndTag() && !isFollowedByClosingTag(tag)) {
      var newTag = new Tag(Tag.SELF_CLOSING_TAG, tagName, tag.getPos(), tag.getLen());
      newTag.setAttributes(tag.getAttributes());
      newTag.setFlags(tag.getFlags());
      tag = newTag
    }
    if (HINT.RULE_TRIM_FIRST_LINE && tag.getFlags() & RULE_TRIM_FIRST_LINE && text[tag.getPos() + tag.getLen()] === "\n") {
      addIgnoreTag(tag.getPos() + tag.getLen(), 1)
    }
    outputTag(tag);
    pushContext(tag);
    createChild(tag)
  }

  function processEndTag(tag) {
    var tagName = tag.getName();
    if (!cntOpen[tagName]) {
      return
    }
    var closeTags = [];
    var i = openTags.length;
    while (--i >= 0) {
      var openTag = openTags[i];
      if (tag.canClose(openTag)) {
        break
      }
      closeTags.push(openTag);
      ++currentFixingCost
    }
    if (i < 0) {
      logger.debug('Skipping end tag with no start tag', {
        'tag': tag
      });
      return
    }
    var flags = tag.getFlags();
    closeTags.forEach(function (openTag) {
      flags |= openTag.getFlags()
    });
    var ignoreWhitespace = (HINT.RULE_IGNORE_WHITESPACE && (flags & RULE_IGNORE_WHITESPACE));
    var keepReopening = HINT.RULE_AUTO_REOPEN && (currentFixingCost < maxFixingCost),
      reopenTags = [];
    closeTags.forEach(function (openTag) {
      var openTagName = openTag.getName();
      if (keepReopening) {
        if (openTag.getFlags() & RULE_AUTO_REOPEN) {
          reopenTags.push(openTag)
        } else {
          keepReopening = !1
        }
      }
      var tagPos = tag.getPos();
      if (ignoreWhitespace) {
        tagPos = getMagicEndPos(tagPos)
      }
      var endTag = new Tag(Tag.END_TAG, openTagName, tagPos, 0);
      endTag.setFlags(openTag.getFlags());
      outputTag(endTag);
      popContext()
    });
    outputTag(tag);
    popContext();
    if (closeTags.length && currentFixingCost < maxFixingCost) {
      var ignorePos = pos;
      i = tagStack.length;
      while (--i >= 0 && ++currentFixingCost < maxFixingCost) {
        var upcomingTag = tagStack[i];
        if (upcomingTag.getPos() > ignorePos || upcomingTag.isStartTag()) {
          break
        }
        var j = closeTags.length;
        while (--j >= 0 && ++currentFixingCost < maxFixingCost) {
          if (upcomingTag.canClose(closeTags[j])) {
            closeTags.splice(j, 1);
            if (reopenTags[j]) {
              reopenTags.splice(j, 1)
            }
            ignorePos = Math.max(ignorePos, upcomingTag.getPos() + upcomingTag.getLen());
            break
          }
        }
      }
      if (ignorePos > pos) {
        outputIgnoreTag(new Tag(Tag.SELF_CLOSING_TAG, 'i', pos, ignorePos - pos))
      }
    }
    reopenTags.forEach(function (startTag) {
      var newTag = addCopyTag(startTag, pos, 0);
      var endTag = startTag.getEndTag();
      if (endTag) {
        newTag.pairWith(endTag)
      }
    })
  }

  function popContext() {
    var tag = openTags.pop();
    --cntOpen[tag.getName()];
    context = context.parentContext
  }

  function pushContext(tag) {
    var tagName = tag.getName(),
      tagFlags = tag.getFlags(),
      tagConfig = tagsConfig[tagName];
    ++cntTotal[tagName];
    if (tag.isSelfClosingTag()) {
      return
    }
    var allowed = [];
    context.allowed.forEach(function (v, k) {
      if (!HINT.RULE_IS_TRANSPARENT || !(tagFlags & RULE_IS_TRANSPARENT)) {
        v = (v & 0xFF00) | (v >> 8)
      }
      allowed.push(tagConfig.allowed[k] & v)
    });
    var flags = tagFlags | (context.flags & RULES_INHERITANCE);
    if (flags & RULE_DISABLE_AUTO_BR) {
      flags &= ~RULE_ENABLE_AUTO_BR
    }
    ++cntOpen[tagName];
    openTags.push(tag);
    context = {
      parentContext: context
    };
    context.allowed = allowed;
    context.flags = flags
  }

  function tagIsAllowed(tagName) {
    var n = tagsConfig[tagName].bitNumber;
    return !!(context.allowed[n >> 3] & (1 << (n & 7)))
  }

  function addStartTag(name, pos, len, prio) {
    return addTag(Tag.START_TAG, name, pos, len, prio || 0)
  }

  function addEndTag(name, pos, len, prio) {
    return addTag(Tag.END_TAG, name, pos, len, prio || 0)
  }

  function addSelfClosingTag(name, pos, len, prio) {
    return addTag(Tag.SELF_CLOSING_TAG, name, pos, len, prio || 0)
  }

  function addBrTag(pos, prio) {
    return addTag(Tag.SELF_CLOSING_TAG, 'br', pos, 0, prio || 0)
  }

  function addIgnoreTag(pos, len, prio) {
    return addTag(Tag.SELF_CLOSING_TAG, 'i', pos, Math.min(len, textLen - pos), prio || 0)
  }

  function addParagraphBreak(pos, prio) {
    return addTag(Tag.SELF_CLOSING_TAG, 'pb', pos, 0, prio || 0)
  }

  function addCopyTag(tag, pos, len, prio) {
    var copy = addTag(tag.getType(), tag.getName(), pos, len, tag.getSortPriority());
    copy.setAttributes(tag.getAttributes());
    return copy
  }

  function addTag(type, name, pos, len, prio) {
    var tag = new Tag(type, name, pos, len, prio || 0);
    if (tagsConfig[name]) {
      tag.setFlags(tagsConfig[name].rules.flags)
    }
    if ((!tagsConfig[name] && !tag.isSystemTag()) || isInvalidTextSpan(pos, len)) {
      tag.invalidate()
    } else if (tagsConfig[name] && tagsConfig[name].isDisabled) {
      logger.warn('Tag is disabled', {
        'tag': tag,
        'tagName': name
      });
      tag.invalidate()
    } else {
      insertTag(tag)
    }
    return tag
  }

  function isInvalidTextSpan(pos, len) {
    return (len < 0 || pos < 0 || pos + len > textLen || /[\uDC00-\uDFFF]/.test(text.substr(pos, 1) + text.substr(pos + len, 1)))
  }

  function insertTag(tag) {
    if (!tagStackIsSorted) {
      tagStack.push(tag)
    } else {
      var i = tagStack.length,
        key = getSortKey(tag);
      while (i > 0 && key > getSortKey(tagStack[i - 1])) {
        tagStack[i] = tagStack[i - 1];
        --i
      }
      tagStack[i] = tag
    }
  }

  function addTagPair(name, startPos, startLen, endPos, endLen, prio) {
    var endTag = addEndTag(name, endPos, endLen, -prio || 0),
      startTag = addStartTag(name, startPos, startLen, prio || 0);
    startTag.pairWith(endTag);
    return startTag
  }

  function addVerbatim(pos, len, prio) {
    return addTag(Tag.SELF_CLOSING_TAG, 'v', pos, len, prio || 0)
  }

  function sortTags() {
    var arr = {},
      keys = [],
      i = tagStack.length;
    while (--i >= 0) {
      var tag = tagStack[i],
        key = getSortKey(tag, i);
      keys.push(key);
      arr[key] = tag
    }
    keys.sort();
    i = keys.length;
    tagStack = [];
    while (--i >= 0) {
      tagStack.push(arr[keys[i]])
    }
    tagStackIsSorted = !0
  }

  function getSortKey(tag, tagIndex) {
    var prioFlag = (tag.getSortPriority() >= 0),
      prio = tag.getSortPriority();
    if (!prioFlag) {
      prio += (1 << 30)
    }
    var lenFlag = (tag.getLen() > 0),
      lenOrder;
    if (lenFlag) {
      lenOrder = textLen - tag.getLen()
    } else {
      var order = {};
      order[Tag.END_TAG] = 0;
      order[Tag.SELF_CLOSING_TAG] = 1;
      order[Tag.START_TAG] = 2;
      lenOrder = order[tag.getType()]
    }
    return hex32(tag.getPos()) + (+prioFlag) + hex32(prio) + (+lenFlag) + hex32(lenOrder) + hex32(tagIndex || 0)
  }

  function hex32(number) {
    var hex = number.toString(16);
    return "        ".substr(hex.length) + hex
  }
  var MSXML = (typeof DOMParser === 'undefined' || typeof XSLTProcessor === 'undefined');
  var xslt = {
    init: function (xsl) {
      var stylesheet = xslt.loadXML(xsl);
      if (MSXML) {
        var generator = new ActiveXObject('MSXML2.XSLTemplate.6.0');
        generator.stylesheet = stylesheet;
        xslt.proc = generator.createProcessor()
      } else {
        xslt.proc = new XSLTProcessor;
        xslt.proc.importStylesheet(stylesheet)
      }
    },
    loadXML: function (xml) {
      var dom;
      if (MSXML) {
        dom = new ActiveXObject('MSXML2.FreeThreadedDOMDocument.6.0');
        dom.async = !1;
        dom.validateOnParse = !1;
        dom.loadXML(xml)
      } else {
        dom = (new DOMParser).parseFromString(xml, 'text/xml')
      }
      if (!dom) {
        throw 'Cannot parse ' + xml
      }
      return dom
    },
    setParameter: function (paramName, paramValue) {
      if (MSXML) {
        xslt.proc.addParameter(paramName, paramValue, '')
      } else {
        xslt.proc.setParameter(null, paramName, paramValue)
      }
    },
    transformToFragment: function (xml, targetDoc) {
      if (MSXML) {
        var div = targetDoc.createElement('div'),
          fragment = targetDoc.createDocumentFragment();
        xslt.proc.input = xslt.loadXML(xml);
        xslt.proc.transform();
        div.innerHTML = xslt.proc.output;
        while (div.firstChild) {
          fragment.appendChild(div.firstChild)
        }
        return fragment
      }
      return xslt.proc.transformToFragment(xslt.loadXML(xml), targetDoc)
    }
  };
  xslt.init(xsl);
  var functionCache = {};

  function preview(text, target) {
    var targetDoc = target.ownerDocument;
    if (!targetDoc) {
      throw 'Target does not have a ownerDocument'
    }
    var resultFragment = xslt.transformToFragment(parse(text).replace(/<[eis]>[^<]*<\/[eis]>/g, ''), targetDoc),
      lastUpdated = target;
    if (typeof window !== 'undefined' && 'chrome' in window) {
      resultFragment.querySelectorAll('script').forEach(function (oldScript) {
        let newScript = document.createElement('script');
        for (let attribute of oldScript.attributes) {
          newScript.setAttribute(attribute.name, attribute.value)
        }
        newScript.textContent = oldScript.textContent;
        oldScript.parentNode.replaceChild(newScript, oldScript)
      })
    }
    if (HINT.hash) {
      computeHashes(resultFragment)
    }
    if (HINT.onRender) {
      executeEvents(resultFragment, 'render')
    }

    function computeHashes(fragment) {
      var nodes = fragment.querySelectorAll('[data-s9e-livepreview-hash]'),
        i = nodes.length;
      while (--i >= 0) {
        nodes[i].setAttribute('data-s9e-livepreview-hash', hash(nodes[i].outerHTML))
      }
    }

    function executeEvent(node, eventName) {
      var code = node.getAttribute('data-s9e-livepreview-on' + eventName);
      if (!functionCache[code]) {
        functionCache[code] = new Function(code)
      }
      functionCache[code].call(node)
    }

    function executeEvents(root, eventName) {
      if (root instanceof Element && root.hasAttribute('data-s9e-livepreview-on' + eventName)) {
        executeEvent(root, eventName)
      }
      var nodes = root.querySelectorAll('[data-s9e-livepreview-on' + eventName + ']'),
        i = nodes.length;
      while (--i >= 0) {
        executeEvent(nodes[i], eventName)
      }
    }

    function refreshElementContent(oldParent, newParent) {
      var oldNodes = oldParent.childNodes,
        newNodes = newParent.childNodes,
        oldCnt = oldNodes.length,
        newCnt = newNodes.length,
        oldNode, newNode, left = 0,
        right = 0;
      while (left < oldCnt && left < newCnt) {
        oldNode = oldNodes[left];
        newNode = newNodes[left];
        if (!refreshNode(oldNode, newNode)) {
          break
        }
        ++left
      }
      var maxRight = Math.min(oldCnt - left, newCnt - left);
      while (right < maxRight) {
        oldNode = oldNodes[oldCnt - (right + 1)];
        newNode = newNodes[newCnt - (right + 1)];
        if (!refreshNode(oldNode, newNode)) {
          break
        }
        ++right
      }
      var i = oldCnt - right;
      while (--i >= left) {
        oldParent.removeChild(oldNodes[i]);
        lastUpdated = oldParent
      }
      var rightBoundary = newCnt - right;
      if (left >= rightBoundary) {
        return
      }
      var newNodesFragment = targetDoc.createDocumentFragment();
      i = left;
      do {
        newNode = newNodes[i];
        if (HINT.onUpdate && newNode instanceof Element) {
          executeEvents(newNode, 'update')
        }
        lastUpdated = newNodesFragment.appendChild(newNode)
      }
      while (i < --rightBoundary);
      if (!right) {
        oldParent.appendChild(newNodesFragment)
      } else {
        oldParent.insertBefore(newNodesFragment, oldParent.childNodes[left])
      }
    }

    function refreshNode(oldNode, newNode) {
      if (oldNode.nodeName !== newNode.nodeName || oldNode.nodeType !== newNode.nodeType) {
        return !1
      }
      if (oldNode instanceof HTMLElement && newNode instanceof HTMLElement) {
        if (!oldNode.isEqualNode(newNode) && !elementHashesMatch(oldNode, newNode)) {
          if (HINT.onUpdate && newNode.hasAttribute('data-s9e-livepreview-onupdate')) {
            executeEvent(newNode, 'update')
          }
          syncElementAttributes(oldNode, newNode);
          refreshElementContent(oldNode, newNode)
        }
      } else if (oldNode.nodeType === 3 || oldNode.nodeType === 8) {
        if (oldNode.nodeValue !== newNode.nodeValue) {
          oldNode.nodeValue = newNode.nodeValue;
          lastUpdated = oldNode
        }
      }
      return !0
    }

    function elementHashesMatch(oldEl, newEl) {
      if (!HINT.hash) {
        return !1
      }
      const attrName = 'data-s9e-livepreview-hash';
      return oldEl.hasAttribute(attrName) && newEl.hasAttribute(attrName) && oldEl.getAttribute(attrName) === newEl.getAttribute(attrName)
    }

    function hash(text) {
      var pos = text.length,
        s1 = 0,
        s2 = 0;
      while (--pos >= 0) {
        s1 = (s1 + text.charCodeAt(pos)) % 0xFFFF;
        s2 = (s1 + s2) % 0xFFFF
      }
      return (s2 << 16) | s1
    }

    function syncElementAttributes(oldEl, newEl) {
      var oldAttributes = oldEl.attributes,
        newAttributes = newEl.attributes,
        oldCnt = oldAttributes.length,
        newCnt = newAttributes.length,
        i = oldCnt,
        ignoreAttrs = ' ' + oldEl.getAttribute('data-s9e-livepreview-ignore-attrs') + ' ';
      while (--i >= 0) {
        var oldAttr = oldAttributes[i],
          namespaceURI = oldAttr.namespaceURI,
          attrName = oldAttr.name;
        if (HINT.ignoreAttrs && ignoreAttrs.indexOf(' ' + attrName + ' ') > -1) {
          continue
        }
        if (!newEl.hasAttributeNS(namespaceURI, attrName)) {
          oldEl.removeAttributeNS(namespaceURI, attrName);
          lastUpdated = oldEl
        }
      }
      i = newCnt;
      while (--i >= 0) {
        var newAttr = newAttributes[i],
          namespaceURI = newAttr.namespaceURI,
          attrName = newAttr.name,
          attrValue = newAttr.value;
        if (HINT.ignoreAttrs && ignoreAttrs.indexOf(' ' + attrName + ' ') > -1) {
          continue
        }
        if (attrValue !== oldEl.getAttributeNS(namespaceURI, attrName)) {
          oldEl.setAttributeNS(namespaceURI, attrName, attrValue);
          lastUpdated = oldEl
        }
      }
    }
    refreshElementContent(target, resultFragment);
    return lastUpdated
  }

  function setParameter(paramName, paramValue) {
    xslt.setParameter(paramName, paramValue)
  }
  if (!window.s9e) window.s9e = {};
  window.s9e.TextFormatter = {
    'preview': preview
  }
})();;


// TODO: Remove this, use https://github.com/flarum/emoji
var module = {};

module.exports = function (e) {
  var d = {};

  function a(u) {
    if (d[u]) return d[u].exports;
    var o = d[u] = {
      i: u,
      l: !1,
      exports: {}
    };
    return e[u].call(o.exports, o, o.exports, a), o.l = !0, o.exports
  }
  return a.m = e, a.c = d, a.d = function (e, d, u) {
    a.o(e, d) || Object.defineProperty(e, d, {
      enumerable: !0,
      get: u
    })
  }, a.r = function (e) {
    "undefined" != typeof Symbol && Symbol.toStringTag && Object.defineProperty(e, Symbol.toStringTag, {
      value: "Module"
    }), Object.defineProperty(e, "__esModule", {
      value: !0
    })
  }, a.t = function (e, d) {
    if (1 & d && (e = a(e)), 8 & d) return e;
    if (4 & d && "object" == typeof e && e && e.__esModule) return e;
    var u = Object.create(null);
    if (a.r(u), Object.defineProperty(u, "default", {
        enumerable: !0,
        value: e
      }), 2 & d && "string" != typeof e)
      for (var o in e) a.d(u, o, function (d) {
        return e[d]
      }.bind(null, o));
    return u
  }, a.n = function (e) {
    var d = e && e.__esModule ? function () {
      return e.default
    } : function () {
      return e
    };
    return a.d(d, "a", d), d
  }, a.o = function (e, d) {
    return Object.prototype.hasOwnProperty.call(e, d)
  }, a.p = "", a(a.s = 10)
}([function (e, d) {
  e.exports = flarum.core.compat.extend
}, function (e) {
  e.exports = JSON.parse('{"😀":["gleeful"],"😃":["glad","smile"],"😄":["happy"],"😁":["blissful","grin"],"😆":["amused","laugh","lol"],"😅":["embarassed"],"🤣":["entertained","rofl"],"😂":["joyful","haha"],"🙂":["pleased"],"🙃":["ecstatic","upside_down"],"😉":["coy","wink"],"😊":["anxious","blush"],"😇":["innocent","halo"],"🥰":["love"],"😍":["lovestruck"],"🤩":["starstruck"],"😘":["flirty"],"😗":["kiss"],"☺️":["relaxed"],"😚":["loving_kiss"],"😙":["happy_kiss"],"😋":["yum","savour"],"😛":["playful","tongue_out"],"😜":["mischievous"],"🤪":["crazy"],"😝":["facetious","lmao"],"🤑":["pretentious","money_mouth"],"🤗":["hugging"],"🤭":["gasp"],"🤫":["shushing"],"🤔":["curious","thinking"],"🤐":["silenced","zipper_mouth"],"🤨":["contempt"],"😐️":["indifferent","neutral"],"😑":["apathetic","expressionless"],"😶":["vacant","no_mouth"],"😏":["cocky","smirk"],"😒":["unamused"],"🙄":["disbelief"],"😬":["grimaced"],"🤥":["lying"],"😌":["relieved"],"😔":["pensive"],"😪":["sleepy"],"🤤":["drooling"],"😴":["exhausted","sleeping"],"😷":["ill","mask"],"🤒":["sick"],"🤕":["injured"],"🤢":["nauseated"],"🤮":["vomiting"],"🤧":["sneezing"],"🥵":["overheating"],"🥶":["freezing"],"🥴":["woozy"],"😵":["dizzy"],"🤯":["shocked","exploding_head"],"🤠":["cowboy"],"🥳":["partying","celebrating"],"😎":["confident"],"🤓":["nerd"],"🧐":["monocle"],"😕":["confused"],"😟":["worried"],"🙁":["cheerless"],"☹️":["sad","frowning"],"😮":["surprised"],"😯":["hushed"],"😲":["astonished"],"😳":["flushed"],"🥺":["pleading"],"😦":["bored"],"😧":["anguished","wtf"],"😨":["fearful"],"😰":["frustrated"],"😥":["hopeful"],"😢":["upset","cry"],"😭":["distressed","sob"],"😱":["frightened","scream"],"😖":["confounded"],"😣":["persevered"],"😞":["disappointed"],"😓":["shamed"],"😩":["weary"],"😫":["tired"],"🥱":["yawn"],"😤":["annoyed","hrmph"],"😡":["enraged","pout"],"😠":["angry"],"🤬":["censored"],"😈":["imp"],"👿":["angry_imp"],"💀":["skull"],"☠️":["crossbones"],"💩":["poop"],"🤡":["clown"],"👹":["ogre"],"👺":["goblin"],"👻":["ghost"],"👽️":["alien"],"👾":["alien_monster","space_invader"],"🤖":["robot"],"😺":["smiling_cat"],"😸":["grinning_cat"],"😹":["joyful_cat"],"😻":["lovestruck_cat"],"😼":["smirking_cat"],"😽":["kissing_cat"],"🙀":["weary_cat"],"😿":["crying_cat"],"😾":["pouting_cat"],"🙈":["see_no_evil"],"🙉":["hear_no_evil"],"🙊":["speak_no_evil"],"💋":["kiss_lips"],"💌":["love_letter"],"💘":["cupid"],"💝":["heart_ribbon"],"💖":["sparkling_heart"],"💗":["heartpulse"],"💓":["heartbeat"],"💞":["revolving_hearts"],"💕":["two_hearts"],"💟":["heart_decoration"],"❣️":["heart_exclamation"],"💔":["broken_heart"],"❤️":["heart"],"🧡":["orange_heart"],"💛":["yellow_heart"],"💚":["green_heart"],"💙":["blue_heart"],"💜":["purple_heart"],"🤎":["brown_heart"],"🖤":["black_heart"],"🤍":["white_heart"],"💯":["100"],"💢":["anger"],"💥":["boom","collision"],"💫":["dizzy_star"],"💦":["sweat_drops"],"💨":["dash"],"🕳️":["hole"],"💣️":["bomb"],"💬":["speech"],"👁‍🗨":["eye_bubble"],"🗨️":["left_speech"],"🗯️":["right_anger_speech"],"💭":["thought"],"💤":["zzz"],"👋":["wave"],"🤚":["raised_backhand"],"🖐️":["splayed_hand"],"✋":["raised_hand"],"🖖":["vulcan"],"👌":["ok_hand"],"🤏":["pinch"],"✌️":["victory"],"🤞":["fingers_crossed"],"🤟":["love_you_gesture"],"🤘":["metal"],"🤙":["call_me"],"👈️":["point_left"],"👉️":["point_right"],"👆️":["backhand_point_up"],"🖕":["middle_finger"],"👇️":["point_down"],"☝️":["point_up"],"👍️":["thumbsup","+1","yes"],"👎️":["thumbsdown","-1","no"],"✊":["fist"],"👊":["punch"],"🤛":["left_facing_fist"],"🤜":["right_facing_fist"],"👏":["clap"],"🙌":["raised_hands"],"👐":["open_hands"],"🤲":["palms_up"],"🤝":["handshake"],"🙏":["pray"],"✍️":["writing_hand"],"💅":["nail_care"],"🤳":["selfie"],"💪":["muscle","right_bicep"],"🦾":["mech_arm"],"🦿":["mech_leg"],"🦵":["leg"],"🦶":["foot"],"👂️":["ear"],"🦻":["hearing_aid"],"👃":["nose"],"🧠":["brain"],"🦷":["tooth"],"🦴":["bone"],"👀":["eyes"],"👁️":["eye"],"👅":["tongue"],"👄":["lips"],"👶":["baby"],"🧒":["child"],"👦":["boy"],"👧":["girl"],"🧑":["adult"],"👱":["blond_person"],"👨":["man"],"🧔":["bearded_person"],"👱‍♂️":["blond_man"],"👨‍🦰":["red_haired_man"],"👨‍🦱":["curly_haired_man"],"👨‍🦳":["white_haired_man"],"👨‍🦲":["bald_man"],"👩":["woman"],"👱‍♀️":["blond_woman"],"👩‍🦰":["red_haired_woman"],"👩‍🦱":["curly_haired_woman"],"👩‍🦳":["white_haired_woman"],"👩‍🦲":["bald_woman"],"🧓":["older_adult"],"👴":["older_man"],"👵":["older_woman"],"🙍":["person_frowning"],"🙍‍♂️":["man_frowning"],"🙍‍♀️":["woman_frowning"],"🙎":["person_pouting"],"🙎‍♂️":["man_pouting"],"🙎‍♀️":["woman_pouting"],"🙅":["person_gesturing_no"],"🙅‍♂️":["man_gesturing_no"],"🙅‍♀️":["woman_gesturing_no"],"🙆":["person_gesturing_ok"],"🙆‍♂️":["man_gesturing_ok"],"🙆‍♀️":["woman_gesturing_ok"],"💁":["person_tipping_hand"],"💁‍♂️":["man_tipping_hand"],"💁‍♀️":["woman_tipping_hand"],"🙋":["person_raising_hand"],"🙋‍♂️":["man_raising_hand"],"🙋‍♀️":["woman_raising_hand"],"🧏":["person_deaf"],"🧏‍♂️":["man_deaf"],"🧏‍♀️":["woman_deaf"],"🙇":["person_bowing"],"🙇‍♂️":["man_bowing"],"🙇‍♀️":["woman_bowing"],"🤦":["person_facepalming"],"🤦‍♂️":["man_facepalming"],"🤦‍♀️":["woman_facepalming"],"🤷":["person_shrugging"],"🤷‍♂️":["man_shrugging"],"🤷‍♀️":["woman_shrugging"],"👨‍⚕️":["man_health_worker"],"👩‍⚕️":["woman_health_worker"],"👨‍🎓":["man_student"],"👩‍🎓":["woman_student"],"👨‍🏫":["man_teacher"],"👩‍🏫":["woman_teacher"],"👨‍⚖️":["man_judge"],"👩‍⚖️":["woman_judge"],"👨‍🌾":["man_farmer"],"👩‍🌾":["woman_farmer"],"👨‍🍳":["man_cook"],"👩‍🍳":["woman_cook"],"👨‍🔧":["man_mechanic"],"👩‍🔧":["woman_mechanic"],"👨‍🏭":["man_factory_worker"],"👩‍🏭":["woman_factory_worker"],"👨‍💼":["man_office_worker"],"👩‍💼":["woman_office_worker"],"👨‍🔬":["man_scientist"],"👩‍🔬":["woman_scientist"],"👨‍💻":["man_technologist"],"👩‍💻":["woman_technologist"],"👨‍🎤":["man_singer"],"👩‍🎤":["woman_singer"],"👨‍🎨":["man_artist"],"👩‍🎨":["woman_artist"],"👨‍✈️":["man_pilot"],"👩‍✈️":["woman_pilot"],"👨‍🚀":["man_astronaut"],"👩‍🚀":["woman_astronaut"],"👨‍🚒":["man_firefighter"],"👩‍🚒":["woman_firefighter"],"👮":["police_officer"],"👮‍♂️":["man_police_officer"],"👮‍♀️":["woman_police_officer"],"🕵️":["detective"],"🕵️‍♂️":["man_detective"],"🕵️‍♀️":["woman_detective"],"💂":["guard"],"💂‍♂️":["man_guard"],"💂‍♀️":["woman_guard"],"👷":["construction_worker"],"👷‍♂️":["man_construction_worker"],"👷‍♀️":["woman_construction_worker"],"🤴":["prince"],"👸":["princess"],"👳":["person_turban"],"👳‍♂️":["man_turban"],"👳‍♀️":["woman_turban"],"👲":["man_chinese_cap"],"🧕":["woman_headscarf"],"🤵":["man_tuxedo"],"👰":["bride_veil"],"🤰":["pregnant_woman"],"🤱":["breast_feeding"],"👼":["baby_angel"],"🎅":["santa"],"🤶":["mrs_claus"],"🦸":["hero","superhero"],"🦸‍♂️":["man_hero"],"🦸‍♀️":["woman_hero"],"🦹":["villain","supervillain"],"🦹‍♂️":["man_villain"],"🦹‍♀️":["woman_villain"],"🧙":["mage"],"🧙‍♂️":["man_mage"],"🧙‍♀️":["woman_mage"],"🧚":["fairy"],"🧚‍♂️":["man_fairy"],"🧚‍♀️":["woman_fairy"],"🧛":["vampire"],"🧛‍♂️":["man_vampire"],"🧛‍♀️":["woman_vampire"],"🧜":["merperson"],"🧜‍♂️":["merman"],"🧜‍♀️":["mermaid"],"🧝":["elf"],"🧝‍♂️":["man_elf"],"🧝‍♀️":["woman_elf"],"🧞":["genie"],"🧞‍♂️":["man_genie"],"🧞‍♀️":["woman_genie"],"🧟":["zombie"],"🧟‍♂️":["man_zombie"],"🧟‍♀️":["woman_zombie"],"💆":["person_getting_massage"],"💆‍♂️":["man_getting_face_massage"],"💆‍♀️":["woman_getting_face_massage"],"💇":["person_getting_haircut"],"💇‍♂️":["man_getting_haircut"],"💇‍♀️":["woman_getting_haircut"],"🚶":["person_walking"],"🚶‍♂️":["man_walking"],"🚶‍♀️":["woman_walking"],"🧍":["person_standing"],"🧍‍♂️":["man_standing"],"🧍‍♀️":["woman_standing"],"🧎":["person_kneeling"],"🧎‍♂️":["man_kneeling"],"🧎‍♀️":["woman_kneeling"],"👨‍🦯":["man_probing_cane"],"👩‍🦯":["woman_probing_cane"],"👨‍🦼":["man_motor_wheelchair"],"👩‍🦼":["woman_motor_wheelchair"],"👨‍🦽":["man_wheelchair"],"👩‍🦽":["woman_wheelchair"],"🏃":["person_running"],"🏃‍♂️":["man_running"],"🏃‍♀️":["woman_running"],"💃":["dancer","woman_dancing"],"🕺":["man_dancing"],"🕴️":["levitate"],"👯":["people_bunny_ears_partying"],"👯‍♂️":["men_bunny_ears_partying"],"👯‍♀️":["women_bunny_ears_partying"],"🧖":["person_steamy_room"],"🧖‍♂️":["man_steamy_room"],"🧖‍♀️":["woman_steamy_room"],"🧗":["person_climbing"],"🧗‍♂️":["man_climbing"],"🧗‍♀️":["woman_climbing"],"🤺":["person_fencing"],"🏇":["horse_racing"],"⛷️":["skier"],"🏂️":["snowboarder"],"🏌️":["person_golfing"],"🏌️‍♂️":["man_golfing"],"🏌️‍♀️":["woman_golfing"],"🏄️":["person_surfing"],"🏄‍♂️":["man_surfing"],"🏄‍♀️":["woman_surfing"],"🚣":["person_rowing_boat"],"🚣‍♂️":["man_rowing_boat"],"🚣‍♀️":["woman_rowing_boat"],"🏊️":["person_swimming"],"🏊‍♂️":["man_swimming"],"🏊‍♀️":["woman_swimming"],"⛹️":["person_bouncing_ball"],"⛹️‍♂️":["man_bouncing_ball"],"⛹️‍♀️":["woman_bouncing_ball"],"🏋️":["person_lifting_weights"],"🏋️‍♂️":["man_lifting_weights"],"🏋️‍♀️":["woman_lifting_weights"],"🚴":["person_biking"],"🚴‍♂️":["man_biking"],"🚴‍♀️":["woman_biking"],"🚵":["person_mountain_biking"],"🚵‍♂️":["man_mountain_biking"],"🚵‍♀️":["woman_mountain_biking"],"🤸":["person_cartwheel"],"🤸‍♂️":["man_cartwheeling"],"🤸‍♀️":["woman_cartwheeling"],"🤼":["people_wrestling"],"🤼‍♂️":["men_wrestling"],"🤼‍♀️":["women_wrestling"],"🤽":["person_water_polo"],"🤽‍♂️":["man_water_polo"],"🤽‍♀️":["woman_water_polo"],"🤾":["person_handball"],"🤾‍♂️":["man_handball"],"🤾‍♀️":["woman_handball"],"🤹":["person_juggling"],"🤹‍♂️":["man_juggling"],"🤹‍♀️":["woman_juggling"],"🧘":["person_lotus_position"],"🧘‍♂️":["man_lotus_position"],"🧘‍♀️":["woman_lotus_position"],"🛀":["bath"],"🛌":["in_bed"],"🧑‍🤝‍🧑":["holding_hands_people"],"👭":["holding_hands_ww"],"👫":["holding_hands_mw","holding_hands_wm"],"👬":["holding_hands_mm"],"💏":["couple"],"👩‍❤️‍💋‍👨":["kiss_mw","kiss_wm"],"👨‍❤️‍💋‍👨":["kiss_mm"],"👩‍❤️‍💋‍👩":["kiss_ww"],"💑":["couple_heart"],"👩‍❤️‍👨":["couple_mw","couple_wm"],"👨‍❤️‍👨":["couple_mm"],"👩‍❤️‍👩":["couple_ww"],"👪️":["family"],"👨‍👩‍👦":["family_mwb"],"👨‍👩‍👧":["family_mwg"],"👨‍👩‍👧‍👦":["family_mwgb"],"👨‍👩‍👦‍👦":["family_mwbb"],"👨‍👩‍👧‍👧":["family_mwgg"],"👨‍👨‍👦":["family_mmb"],"👨‍👨‍👧":["family_mmg"],"👨‍👨‍👧‍👦":["family_mmgb"],"👨‍👨‍👦‍👦":["family_mmbb"],"👨‍👨‍👧‍👧":["family_mmgg"],"👩‍👩‍👦":["family_wwb"],"👩‍👩‍👧":["family_wwg"],"👩‍👩‍👧‍👦":["family_wwgb"],"👩‍👩‍👦‍👦":["family_wwbb"],"👩‍👩‍👧‍👧":["family_wwgg"],"👨‍👦":["family_mb"],"👨‍👦‍👦":["family_mbb"],"👨‍👧":["family_mg"],"👨‍👧‍👦":["family_mgb"],"👨‍👧‍👧":["family_mgg"],"👩‍👦":["family_wb"],"👩‍👦‍👦":["family_wbb"],"👩‍👧":["family_wg"],"👩‍👧‍👦":["family_wgb"],"👩‍👧‍👧":["family_wgg"],"🗣️":["speaking_head"],"👤":["bust_silhouette"],"👥":["busts_silhouette"],"👣":["footprints"],"🏻":["tone_light","tone1"],"🏼":["tone_medium_light","tone2"],"🏽":["tone_medium","tone3"],"🏾":["tone_medium_dark","tone4"],"🏿":["tone_dark","tone5"],"🦰":["red_hair"],"🦱":["curly_hair"],"🦳":["white_hair"],"🦲":["bald"],"🐵":["monkey_face"],"🐒":["monkey"],"🦍":["gorilla"],"🦧":["orangutan"],"🐶":["dog_face"],"🐕️":["dog"],"🦮":["guide_dog"],"🐕‍🦺":["service_dog"],"🐩":["poodle"],"🐺":["wolf_face"],"🦊":["fox_face"],"🦝":["raccoon"],"🐱":["cat_face"],"🐈️":["cat"],"🦁":["lion_face"],"🐯":["tiger_face"],"🐅":["tiger"],"🐆":["leopard"],"🐴":["horse_face"],"🐎":["horse"],"🦄":["unicorn_face"],"🦓":["zebra"],"🦌":["deer"],"🐮":["cow_face"],"🐂":["ox"],"🐃":["water_buffalo"],"🐄":["cow"],"🐷":["pig_face"],"🐖":["pig"],"🐗":["boar"],"🐽":["pig_nose"],"🐏":["ram"],"🐑":["sheep"],"🐐":["goat"],"🐪":["camel"],"🐫":["two_hump_camel"],"🦙":["llama"],"🦒":["giraffe"],"🐘":["elephant"],"🦏":["rhino"],"🦛":["hippo"],"🐭":["mouse_face"],"🐁":["mouse"],"🐀":["rat"],"🐹":["hamster_face"],"🐰":["rabbit_face"],"🐇":["rabbit"],"🐿️":["chipmunk"],"🦔":["hedgehog"],"🦇":["bat"],"🐻":["bear_face"],"🐨":["koala_face"],"🐼":["panda_face"],"🦥":["sloth"],"🦦":["otter"],"🦨":["skunk"],"🦘":["kangaroo"],"🦡":["badger"],"🐾":["feet"],"🦃":["turkey"],"🐔":["chicken"],"🐓":["rooster"],"🐣":["hatching_chick"],"🐤":["baby_chick"],"🐥":["hatched_chick"],"🐦️":["bird"],"🐧":["penguin"],"🕊️":["dove"],"🦅":["eagle"],"🦆":["duck"],"🦢":["swan"],"🦉":["owl"],"🦩":["flamingo"],"🦚":["peacock"],"🦜":["parrot"],"🐸":["frog_face"],"🐊":["crocodile"],"🐢":["turtle"],"🦎":["lizard"],"🐍":["snake"],"🐲":["dragon_face"],"🐉":["dragon"],"🦕":["sauropod"],"🦖":["trex"],"🐳":["spouting_whale"],"🐋":["whale"],"🐬":["dolphin"],"🐟️":["fish"],"🐠":["tropical_fish"],"🐡":["blowfish"],"🦈":["shark"],"🐙":["octopus"],"🐚":["shell"],"🐌":["snail"],"🦋":["butterfly"],"🐛":["bug"],"🐜":["ant"],"🐝":["bee"],"🐞":["beetle"],"🦗":["cricket"],"🕷️":["spider"],"🕸️":["spider_web"],"🦂":["scorpion"],"🦟":["mosquito"],"🦠":["microbe","germ"],"💐":["bouquet"],"🌸":["cherry_blossom"],"💮":["white_flower"],"🏵️":["rosette"],"🌹":["rose"],"🥀":["wilted_rose"],"🌺":["hibiscus"],"🌻":["sunflower"],"🌼":["blossom"],"🌷":["tulip"],"🌱":["seedling"],"🌲":["evergreen_tree"],"🌳":["deciduous_tree"],"🌴":["palm_tree"],"🌵":["cactus"],"🌾":["ear_of_rice"],"🌿":["herb"],"☘️":["shamrock"],"🍀":["four_leaf_clover"],"🍁":["maple_leaf"],"🍂":["fallen_leaf"],"🍃":["leaves"],"🍇":["grapes"],"🍈":["melon"],"🍉":["watermelon"],"🍊":["tangerine"],"🍋":["lemon"],"🍌":["banana"],"🍍":["pineapple"],"🥭":["mango"],"🍎":["apple"],"🍏":["green_apple"],"🍐":["pear"],"🍑":["peach"],"🍒":["cherries"],"🍓":["strawberry"],"🥝":["kiwi"],"🍅":["tomato"],"🥥":["coconut"],"🥑":["avocado"],"🍆":["eggplant"],"🥔":["potato"],"🥕":["carrot"],"🌽":["corn"],"🌶️":["hot_pepper"],"🥒":["cucumber"],"🥬":["leafy_green"],"🥦":["broccoli"],"🧄":["garlic"],"🧅":["onion"],"🍄":["mushroom"],"🥜":["peanuts"],"🌰":["chestnut"],"🍞":["bread"],"🥐":["croissant"],"🥖":["french_bread"],"🥨":["pretzel"],"🥯":["bagel"],"🥞":["pancakes"],"🧇":["waffle"],"🧀":["cheese"],"🍖":["meat_on_bone"],"🍗":["poultry_leg"],"🥩":["cut_of_meat"],"🥓":["bacon"],"🍔":["hamburger"],"🍟":["fries"],"🍕":["pizza"],"🌭":["hotdog"],"🥪":["sandwich"],"🌮":["taco"],"🌯":["burrito"],"🥙":["stuffed_flatbread"],"🧆":["falafel"],"🥚":["egg"],"🍳":["cooking"],"🥘":["shallow_pan_of_food"],"🍲":["stew"],"🥣":["bowl_spoon"],"🥗":["salad"],"🍿":["popcorn"],"🧈":["butter"],"🧂":["salt"],"🥫":["canned_food"],"🍱":["bento"],"🍘":["rice_cracker"],"🍙":["rice_ball"],"🍚":["rice"],"🍛":["curry"],"🍜":["ramen"],"🍝":["spaghetti"],"🍠":["sweet_potato"],"🍢":["oden"],"🍣":["sushi"],"🍤":["fried_shrimp"],"🍥":["fish_cake"],"🥮":["moon_cake"],"🍡":["dango"],"🥟":["dumpling"],"🥠":["fortune_cookie"],"🥡":["takeout_box"],"🦀":["crab"],"🦞":["lobster"],"🦐":["shrimp"],"🦑":["squid"],"🦪":["oyster"],"🍦":["icecream"],"🍧":["shaved_ice"],"🍨":["ice_cream"],"🍩":["doughnut"],"🍪":["cookie"],"🎂":["birthday"],"🍰":["cake"],"🧁":["cupcake"],"🥧":["pie"],"🍫":["chocolate_bar"],"🍬":["candy"],"🍭":["lollipop"],"🍮":["custard"],"🍯":["honey_pot"],"🍼":["baby_bottle"],"🥛":["milk"],"☕️":["coffee"],"🍵":["tea"],"🍶":["sake"],"🍾":["champagne"],"🍷":["wine_glass"],"🍸️":["cocktail"],"🍹":["tropical_drink"],"🍺":["beer"],"🍻":["beers"],"🥂":["champagne_glass"],"🥃":["tumbler_glass"],"🥤":["cup_straw"],"🧃":["beverage_box","juice_box"],"🧉":["mate","yerba_mate"],"🧊":["ice"],"🥢":["chopsticks"],"🍽️":["fork_knife_plate"],"🍴":["utensils"],"🥄":["spoon"],"🔪":["knife"],"🏺":["amphora"],"🌍️":["earth_africa"],"🌎️":["earth_americas"],"🌏️":["earth_asia"],"🌐":["globe"],"🗺️":["map"],"🗾":["japan"],"🧭":["compass"],"🏔️":["snowy_mountain"],"⛰️":["mountain"],"🌋":["volcano"],"🗻":["mount_fuji"],"🏕️":["camping"],"🏖️":["beach"],"🏜️":["desert"],"🏝️":["island"],"🏞️":["park"],"🏟️":["stadium"],"🏛️":["classical_building"],"🏗️":["construction_site"],"🧱":["brick"],"🏘️":["homes"],"🏚️":["house_abandoned"],"🏠️":["house"],"🏡":["house_garden"],"🏢":["office"],"🏣":["ja_post_office"],"🏤":["post_office"],"🏥":["hospital"],"🏦":["bank"],"🏨":["hotel"],"🏩":["love_hotel"],"🏪":["convenience_store"],"🏫":["school"],"🏬":["department_store"],"🏭️":["factory"],"🏯":["japanese_castle"],"🏰":["castle","european_castle"],"💒":["wedding"],"🗼":["tokyo_tower"],"🗽":["statue_of_liberty"],"⛪️":["church"],"🕌":["mosque"],"🛕":["hindu_temple"],"🕍":["synagogue"],"⛩️":["shinto_shrine"],"🕋":["kaaba"],"⛲️":["fountain"],"⛺️":["tent"],"🌁":["foggy"],"🌃":["night_stars"],"🏙️":["cityscape"],"🌄":["sunrise_over_mountains"],"🌅":["sunrise"],"🌆":["dusk"],"🌇":["sunset"],"🌉":["bridge_at_night"],"♨️":["hotsprings"],"🎠":["carousel_horse"],"🎡":["ferris_wheel"],"🎢":["roller_coaster"],"💈":["barber"],"🎪":["circus_tent"],"🚂":["steam_locomotive"],"🚃":["railway_car"],"🚄":["bullettrain_side"],"🚅":["bullettrain"],"🚆":["train"],"🚇️":["metro"],"🚈":["light_rail"],"🚉":["station"],"🚊":["tram"],"🚝":["monorail"],"🚞":["mountain_railway"],"🚋":["tram_car"],"🚌":["bus"],"🚍️":["oncoming_bus"],"🚎":["trolleybus"],"🚐":["minibus"],"🚑️":["ambulance"],"🚒":["fire_engine"],"🚓":["police_car"],"🚔️":["oncoming_police_car"],"🚕":["taxi"],"🚖":["oncoming_taxi"],"🚗":["red_car"],"🚘️":["oncoming_automobile"],"🚙":["blue_car"],"🚚":["truck"],"🚛":["lorry"],"🚜":["tractor"],"🏎️":["race_car"],"🏍️":["motorcycle"],"🛵":["motor_scooter"],"🦽":["wheelchair"],"🦼":["motor_wheelchair"],"🛺":["auto_rickshaw"],"🚲️":["bike"],"🛴":["scooter"],"🛹":["skateboard"],"🚏":["bus_stop"],"🛣️":["motorway"],"🛤️":["railway_track"],"🛢️":["oil_drum"],"⛽️":["fuel_pump"],"🚨":["rotating_light","police_light"],"🚥":["traffic_light"],"🚦":["vertical_traffic_light"],"🛑":["stop_sign","octagonal_sign"],"🚧":["construction"],"⚓️":["anchor"],"⛵️":["sailboat"],"🛶":["canoe"],"🚤":["speedboat"],"🛳️":["cruise_ship"],"⛴️":["ferry"],"🛥️":["motorboat"],"🚢":["ship"],"✈️":["airplane"],"🛩️":["small_airplane"],"🛫":["airplane_departure"],"🛬":["airplane_arriving"],"🪂":["parachute"],"💺":["seat"],"🚁":["helicopter"],"🚟":["suspension_railway"],"🚠":["mountain_cableway"],"🚡":["aerial_tramway"],"🛰️":["satellite"],"🚀":["rocket"],"🛸":["flying_saucer"],"🛎️":["bellhop"],"🧳":["luggage"],"⌛️":["hourglass"],"⏳️":["hourglass_flowing"],"⌚️":["watch"],"⏰":["alarm_clock"],"⏱️":["stopwatch"],"⏲️":["timer"],"🕰️":["clock"],"🕛️":["clock12"],"🕧️":["clock1230"],"🕐️":["clock1"],"🕜️":["clock130"],"🕑️":["clock2"],"🕝️":["clock230"],"🕒️":["clock3"],"🕞️":["clock330"],"🕓️":["clock4"],"🕟️":["clock430"],"🕔️":["clock5"],"🕠️":["clock530"],"🕕️":["clock6"],"🕡️":["clock630"],"🕖️":["clock7"],"🕢️":["clock730"],"🕗️":["clock8"],"🕣️":["clock830"],"🕘️":["clock9"],"🕤️":["clock930"],"🕙️":["clock10"],"🕥️":["clock1030"],"🕚️":["clock11"],"🕦️":["clock1130"],"🌑":["new_moon"],"🌒":["waxing_crescent_moon"],"🌓":["first_quarter_moon"],"🌔":["waxing_gibbous_moon"],"🌕️":["full_moon"],"🌖":["waning_gibbous_moon"],"🌗":["last_quarter_moon"],"🌘":["waning_crescent_moon"],"🌙":["crescent_moon"],"🌚":["new_moon_face"],"🌛":["first_quarter_moon_face"],"🌜️":["last_quarter_moon_face"],"🌡️":["thermometer"],"☀️":["sun"],"🌝":["full_moon_face"],"🌞":["sun_face"],"🪐":["ringed_planet","saturn"],"⭐️":["star"],"🌟":["star2","glowing_star"],"🌠":["star3","shooting_star"],"🌌":["milky_way"],"☁️":["cloud"],"⛅️":["partly_sunny"],"⛈️":["storm"],"🌤️":["overcast"],"🌥️":["cloudy"],"🌦️":["sunshower"],"🌧️":["rain"],"🌨️":["snow"],"🌩️":["lightning"],"🌪️":["tornado"],"🌫️":["fog"],"🌬️":["wind_face"],"🌀":["cyclone"],"🌈":["rainbow"],"🌂":["closed_umbrella"],"☂️":["umbrella"],"☔️":["umbrella_rain"],"⛱️":["beach_umbrella"],"⚡️":["zap","high_voltage"],"❄️":["snowflake"],"☃️":["snowy_snowman"],"⛄️":["snowman"],"☄️":["comet"],"🔥":["fire"],"💧":["droplet"],"🌊":["ocean"],"🎃":["jack_o_lantern"],"🎄":["christmas_tree","xmas_tree"],"🎆":["fireworks"],"🎇":["sparkler"],"🧨":["firecracker"],"✨":["sparkles"],"🎈":["balloon"],"🎉":["tada","party"],"🎊":["confetti_ball"],"🎋":["tanabata_tree"],"🎍":["bamboo","pine_decor"],"🎎":["dolls"],"🎏":["carp_streamer"],"🎐":["wind_chime"],"🎑":["moon_ceremony","rice_scene"],"🧧":["red_envelope"],"🎀":["ribbon"],"🎁":["gift"],"🎗️":["reminder_ribbon"],"🎟️":["tickets","admission"],"🎫":["ticket"],"🎖️":["military_medal"],"🏆️":["trophy"],"🏅":["medal"],"🥇":["first_place"],"🥈":["second_place"],"🥉":["third_place"],"⚽️":["soccer"],"⚾️":["baseball"],"🥎":["softball"],"🏀":["basketball"],"🏐":["volleyball"],"🏈":["football"],"🏉":["rugby"],"🎾":["tennis"],"🥏":["flying_disc"],"🎳":["bowling"],"🏏":["cricket_game"],"🏑":["field_hockey"],"🏒":["hockey"],"🥍":["lacrosse"],"🏓":["ping_pong"],"🏸":["badminton"],"🥊":["boxing_glove"],"🥋":["gi","martial_arts_uniform"],"🥅":["goal"],"⛳️":["golf"],"⛸️":["ice_skate"],"🎣":["fishing_pole"],"🤿":["diving_mask","scuba_mask"],"🎽":["running_shirt"],"🎿":["ski"],"🛷":["sled"],"🥌":["curling_stone"],"🎯":["dart"],"🪀":["yoyo"],"🪁":["kite"],"🎱":["8ball"],"🔮":["crystal_ball"],"🧿":["nazar_amulet"],"🎮️":["video_game"],"🕹️":["joystick"],"🎰":["slot_machine"],"🎲":["game_die"],"🧩":["jigsaw","puzzle_piece"],"🧸":["teddy_bear"],"♠️":["spades"],"♥️":["hearts"],"♦️":["diamonds"],"♣️":["clubs"],"♟️":["chess_pawn"],"🃏":["black_joker"],"🀄️":["mahjong"],"🎴":["flower_cards"],"🎭️":["performing_arts"],"🖼️":["frame_photo"],"🎨":["art","palette"],"🧵":["spool"],"🧶":["yarn"],"👓️":["glasses"],"🕶️":["sunglasses"],"🥽":["goggles"],"🥼":["lab_coat"],"🦺":["safety_vest"],"👔":["necktie","tie"],"👕":["shirt"],"👖":["jeans"],"🧣":["scarf"],"🧤":["gloves"],"🧥":["coat"],"🧦":["socks"],"👗":["dress"],"👘":["kimono"],"🥻":["sari"],"🩱":["one_piece_swimsuit"],"🩲":["briefs"],"🩳":["shorts"],"👙":["bikini"],"👚":["blouse","womans_clothes"],"👛":["purse"],"👜":["handbag"],"👝":["pouch","clutch_bag"],"🛍️":["shopping_bags"],"🎒":["backpack"],"👞":["dress_shoe","mans_shoe"],"👟":["sneaker","athletic_shoe"],"🥾":["hiking_boot"],"🥿":["flat_shoe"],"👠":["high_heel"],"👡":["womans_sandal"],"🩰":["ballet_shoes"],"👢":["womans_boot"],"👑":["crown"],"👒":["womans_hat"],"🎩":["top_hat"],"🎓️":["graduation_cap"],"🧢":["billed_cap"],"⛑️":["helmet_cross"],"📿":["prayer_beads"],"💄":["lipstick"],"💍":["ring"],"💎":["gem"],"🔇":["mute","no_sound"],"🔈️":["speaker","low_sound"],"🔉":["sound"],"🔊":["loud_sound"],"📢":["loudspeaker"],"📣":["megaphone"],"📯":["postal_horn"],"🔔":["bell"],"🔕":["no_bell"],"🎼":["musical_score"],"🎵":["musical_note"],"🎶":["musical_notes"],"🎙️":["studio_microphone"],"🎚️":["level_slider"],"🎛️":["control_knobs"],"🎤":["microphone"],"🎧️":["headphones"],"📻️":["radio"],"🎷":["saxophone"],"🎸":["guitar"],"🎹":["musical_keyboard"],"🎺":["trumpet"],"🎻":["violin"],"🪕":["banjo"],"🥁":["drum"],"📱":["mobile","iphone","android"],"📲":["mobile_calling"],"☎️":["telephone"],"📞":["telephone_receiver"],"📟️":["pager"],"📠":["fax"],"🔋":["battery"],"🔌":["electric_plug"],"💻️":["laptop"],"🖥️":["desktop","computer"],"🖨️":["printer"],"⌨️":["keyboard"],"🖱️":["computer_mouse"],"🖲️":["trackball"],"💽":["minidisc"],"💾":["floppy_disk"],"💿️":["cd","disk"],"📀":["dvd"],"🧮":["abacus"],"🎥":["movie_camera"],"🎞️":["film_frames"],"📽️":["projector"],"🎬️":["clapper"],"📺️":["tv"],"📷️":["camera"],"📸":["camera_flash"],"📹️":["video_camera"],"📼":["vhs"],"🔍️":["mag"],"🔎":["mag_right"],"🕯️":["candle"],"💡":["bulb","light_bulb"],"🔦":["flashlight"],"🏮":["red_lantern"],"🪔":["diya_lamp"],"📔":["decorative_notebook"],"📕":["closed_book"],"📖":["book"],"📗":["green_book"],"📘":["blue_book"],"📙":["orange_book"],"📚️":["books"],"📓":["notebook"],"📒":["ledger"],"📃":["page_curl"],"📜":["scroll"],"📄":["page_facing_up"],"📰":["newspaper"],"🗞️":["rolled_newspaper"],"📑":["bookmark_tabs"],"🔖":["bookmark"],"🏷️":["label"],"💰️":["moneybag"],"💴":["yen"],"💵":["dollar"],"💶":["euro"],"💷":["pound"],"💸":["money_wings"],"💳️":["credit_card"],"🧾":["receipt"],"💹":["ja_chart"],"💱":["currency_exchange"],"💲":["dollar_sign"],"✉️":["envelope"],"📧":["email"],"📨":["incoming_envelope"],"📩":["envelope_arrow"],"📤️":["outbox_tray"],"📥️":["inbox_tray"],"📦️":["package"],"📫️":["mailbox"],"📪️":["mailbox_closed"],"📬️":["mailbox_mail"],"📭️":["mailbox_no_mail"],"📮":["postbox"],"🗳️":["ballot_box"],"✏️":["pencil"],"✒️":["black_nib"],"🖋️":["fountain_pen"],"🖊️":["pen"],"🖌️":["paintbrush"],"🖍️":["crayon"],"📝":["memo"],"💼":["briefcase"],"📁":["file_folder"],"📂":["open_file_folder"],"🗂️":["dividers"],"📅":["date","calendar"],"📆":["torn_calendar"],"🗒️":["notepad_spiral"],"🗓️":["calendar_spiral"],"📇":["card_index"],"📈":["chart_up"],"📉":["chart_down"],"📊":["bar_chart"],"📋️":["clipboard"],"📌":["pushpin"],"📍":["round_pushpin"],"📎":["paperclip"],"🖇️":["paperclips"],"📏":["straight_ruler"],"📐":["triangular_ruler"],"✂️":["scissors"],"🗃️":["card_box"],"🗄️":["file_cabinet"],"🗑️":["trashcan","wastebasket"],"🔒️":["lock"],"🔓️":["unlock"],"🔏":["locked_pen"],"🔐":["locked_key"],"🔑":["key"],"🗝️":["old_key"],"🔨":["hammer"],"🪓":["axe"],"⛏️":["pick"],"⚒️":["hammer_pick"],"🛠️":["tools","hammer_wrench"],"🗡️":["dagger"],"⚔️":["crossed_swords"],"🔫":["gun","pistol"],"🏹":["bow"],"🛡️":["shield"],"🔧":["wrench"],"🔩":["nut_and_bolt"],"⚙️":["gear"],"🗜️":["clamp","compression"],"⚖️":["scales"],"🦯":["probing_cane"],"🔗":["link"],"⛓️":["chains"],"🧰":["toolbox"],"🧲":["magnet"],"⚗️":["alembic"],"🧪":["test_tube"],"🧫":["petri_dish"],"🧬":["dna","double_helix"],"🔬":["microscope"],"🔭":["telescope"],"📡":["satellite_antenna"],"💉":["syringe"],"🩸":["blood_drop"],"💊":["pill"],"🩹":["bandaid","adhesive_bandage"],"🩺":["stethoscope"],"🚪":["door"],"🛏️":["bed"],"🛋️":["couch"],"🪑":["chair"],"🚽":["toilet"],"🚿":["shower"],"🛁":["bathtub"],"🪒":["razor"],"🧴":["lotion"],"🧷":["safety_pin"],"🧹":["broom"],"🧺":["basket"],"🧻":["toilet_paper"],"🧼":["soap"],"🧽":["sponge"],"🧯":["fire_extinguisher"],"🛒":["shopping_cart"],"🚬":["cigarette","smoking"],"⚰️":["coffin"],"⚱️":["urn"],"🗿":["moai"],"🏧":["atm"],"🚮":["litter_bin"],"🚰":["potable_water"],"♿️":["handicapped"],"🚹️":["mens"],"🚺️":["womens"],"🚻":["restroom","bathroom"],"🚼️":["baby_symbol"],"🚾":["wc"],"🛂":["passport_control"],"🛃":["customs"],"🛄":["baggage_claim"],"🛅":["left_luggage"],"⚠️":["warning"],"🚸":["children_crossing"],"⛔️":["no_entry"],"🚫":["no_entry_sign"],"🚳":["no_bicycles"],"🚭️":["no_smoking"],"🚯":["do_not_litter"],"🚱":["non_potable_water"],"🚷":["no_pedestrians"],"📵":["no_mobile_phones"],"🔞":["underage"],"☢️":["radioactive"],"☣️":["biohazard"],"⬆️":["arrow_up"],"↗️":["arrow_upper_right"],"➡️":["arrow_right"],"↘️":["arrow_lower_right"],"⬇️":["arrow_down"],"↙️":["arrow_lower_left"],"⬅️":["arrow_left"],"↖️":["arrow_upper_left"],"↕️":["arrow_up_down"],"↔️":["arrow_left_right"],"↩️":["arrow_left_hook"],"↪️":["arrow_right_hook"],"⤴️":["arrow_heading_up"],"⤵️":["arrow_heading_down"],"🔃":["clockwise"],"🔄":["counter_clockwise"],"🔙":["back"],"🔚":["end"],"🔛":["on"],"🔜":["soon"],"🔝":["top"],"🛐":["place_of_worship"],"⚛️":["atom"],"🕉️":["om_symbol"],"✡️":["star_of_david"],"☸️":["wheel_of_dharma"],"☯️":["yin_yang"],"✝️":["cross"],"☦️":["orthodox_cross"],"☪️":["star_and_crescent"],"☮️":["peace"],"🕎":["menorah"],"🔯":["six_pointed_star"],"♈️":["aries"],"♉️":["taurus"],"♊️":["gemini"],"♋️":["cancer"],"♌️":["leo"],"♍️":["virgo"],"♎️":["libra"],"♏️":["scorpius"],"♐️":["sagittarius"],"♑️":["capricorn"],"♒️":["aquarius"],"♓️":["pisces"],"⛎":["ophiuchus"],"🔀":["shuffle"],"🔁":["repeat"],"🔂":["repeat_single"],"▶️":["play"],"⏩️":["fast_forward"],"⏭️":["next_track"],"⏯️":["play_pause"],"◀️":["reverse"],"⏪️":["rewind"],"⏮️":["previous_track"],"🔼":["up_button"],"⏫":["fast_up_button"],"🔽":["down_button"],"⏬":["fast_down_button"],"⏸️":["pause"],"⏹️":["stop"],"⏺️":["record"],"⏏️":["eject"],"🎦":["cinema"],"🔅":["dim","low_brightness"],"🔆":["bright","high_brightness"],"📶":["signal_strength","antenna_bars"],"📳":["vibration_mode"],"📴":["mobile_phone_off"],"♀️":["female","female_sign"],"♂️":["male","male_sign"],"⚕️":["medical"],"♾️":["infinity"],"♻️":["recycle"],"⚜️":["fleur-de-lis"],"🔱":["trident"],"📛":["name_badge"],"🔰":["ja_beginner"],"⭕️":["o"],"✅":["white_check_mark"],"☑️":["checked_ballot"],"✔️":["check_mark"],"✖️":["multiplication"],"❌":["x","cross_mark"],"❎":["cross_mark_button"],"➕":["plus"],"➖":["minus"],"➗":["division"],"➰":["curly_loop"],"➿":["double_curly_loop"],"〽️":["part_alternation_mark"],"✳️":["eight_spoked_asterisk"],"✴️":["eight_pointed_star"],"❇️":["sparkle"],"‼️":["bangbang","double_exclamation"],"⁉️":["interrobang","exclamation_question"],"❓️":["question"],"❔":["white_question"],"❕":["white_exclamation"],"❗️":["exclamation"],"〰️":["wavy_dash"],"©️":["copyright"],"®️":["registered"],"™️":["tm"],"#️⃣":["hash"],"*️⃣":["asterisk"],"0️⃣":["zero"],"1️⃣":["one"],"2️⃣":["two"],"3️⃣":["three"],"4️⃣":["four"],"5️⃣":["five"],"6️⃣":["six"],"7️⃣":["seven"],"8️⃣":["eight"],"9️⃣":["nine"],"🔟":["ten"],"🔠":["upper_abcd"],"🔡":["abcd"],"🔢":["1234"],"🔣":["symbols"],"🔤":["abc"],"🅰️":["a_blood"],"🆎":["ab_blood"],"🅱️":["b_blood"],"🆑":["cl"],"🆒":["cool"],"🆓":["free"],"ℹ️":["info"],"🆔":["id"],"Ⓜ️":["m"],"🆕":["new"],"🆖":["ng"],"🅾️":["o_blood"],"🆗":["ok"],"🅿️":["p"],"🆘":["sos"],"🆙":["up"],"🆚":["vs"],"🈁":["ja_here","koko"],"🈂️":["ja_service_charge"],"🈷️":["ja_monthly_amount"],"🈶":["ja_not_free_of_carge"],"🈯️":["ja_reserved"],"🉐":["ja_bargain"],"🈹":["ja_discount"],"🈚️":["ja_free_of_charge"],"🈲":["ja_prohibited"],"🉑":["ja_acceptable"],"🈸":["ja_application"],"🈴":["ja_passing_grade"],"🈳":["ja_vacancy"],"㊗️":["ja_congratulations"],"㊙️":["ja_secret"],"🈺":["ja_open_for_business"],"🈵":["ja_no_vacancy"],"🔴":["red_circle"],"🟠":["orange_circle"],"🟡":["yellow_circle"],"🟢":["green_circle"],"🔵":["blue_circle"],"🟣":["purple_circle"],"🟤":["brown_circle"],"⚫️":["black_circle"],"⚪️":["white_circle"],"🟥":["red_square"],"🟧":["orange_square"],"🟨":["yellow_square"],"🟩":["green_square"],"🟦":["blue_square"],"🟪":["purple_square"],"🟫":["brown_square"],"⬛️":["large_black_square"],"⬜️":["large_white_square"],"◼️":["medium_black_square"],"◻️":["medium_white_square"],"◾️":["medium_small_black_square"],"◽️":["medium_small_white_square"],"▪️":["small_black_square"],"▫️":["small_white_square"],"🔶":["large_orange_diamond"],"🔷":["large_blue_diamond"],"🔸":["small_orange_diamond"],"🔹":["small_blue_diamond"],"🔺":["up_red_triangle"],"🔻":["down_red_triangle"],"💠":["diamond_dot"],"🔘":["radio_button"],"🔳":["white_square_button"],"🔲":["black_square_button"],"🏁":["checkered_flag"],"🚩":["triangle_flag"],"🎌":["crossed_flags"],"🏴":["black_flag"],"🏳️":["white_flag"],"🏳️‍🌈":["rainbow_flag"],"🏴‍☠️":["pirate_flag","jolly_roger"],"🇦🇨":["flag_ac"],"🇦🇩":["flag_ad"],"🇦🇪":["flag_ae"],"🇦🇫":["flag_af"],"🇦🇬":["flag_ag"],"🇦🇮":["flag_ai"],"🇦🇱":["flag_al"],"🇦🇲":["flag_am"],"🇦🇴":["flag_ao"],"🇦🇶":["flag_aq"],"🇦🇷":["flag_ar"],"🇦🇸":["flag_as"],"🇦🇹":["flag_at"],"🇦🇺":["flag_au"],"🇦🇼":["flag_aw"],"🇦🇽":["flag_ax"],"🇦🇿":["flag_az"],"🇧🇦":["flag_ba"],"🇧🇧":["flag_bb"],"🇧🇩":["flag_bd"],"🇧🇪":["flag_be"],"🇧🇫":["flag_bf"],"🇧🇬":["flag_bg"],"🇧🇭":["flag_bh"],"🇧🇮":["flag_bi"],"🇧🇯":["flag_bj"],"🇧🇱":["flag_bl"],"🇧🇲":["flag_bm"],"🇧🇳":["flag_bn"],"🇧🇴":["flag_bo"],"🇧🇶":["flag_bq"],"🇧🇷":["flag_br"],"🇧🇸":["flag_bs"],"🇧🇹":["flag_bt"],"🇧🇻":["flag_bv"],"🇧🇼":["flag_bw"],"🇧🇾":["flag_by"],"🇧🇿":["flag_bz"],"🇨🇦":["flag_ca"],"🇨🇨":["flag_cc"],"🇨🇩":["flag_cd"],"🇨🇫":["flag_cf"],"🇨🇬":["flag_cg"],"🇨🇭":["flag_ch"],"🇨🇮":["flag_ci"],"🇨🇰":["flag_ck"],"🇨🇱":["flag_cl"],"🇨🇲":["flag_cm"],"🇨🇳":["flag_cn"],"🇨🇴":["flag_co"],"🇨🇵":["flag_cp"],"🇨🇷":["flag_cr"],"🇨🇺":["flag_cu"],"🇨🇻":["flag_cv"],"🇨🇼":["flag_cw"],"🇨🇽":["flag_cx"],"🇨🇾":["flag_cy"],"🇨🇿":["flag_cz"],"🇩🇪":["flag_de"],"🇩🇬":["flag_dg"],"🇩🇯":["flag_dj"],"🇩🇰":["flag_dk"],"🇩🇲":["flag_dm"],"🇩🇴":["flag_do"],"🇩🇿":["flag_dz"],"🇪🇦":["flag_ea"],"🇪🇨":["flag_ec"],"🇪🇪":["flag_ee"],"🇪🇬":["flag_eg"],"🇪🇭":["flag_eh"],"🇪🇷":["flag_er"],"🇪🇸":["flag_es"],"🇪🇹":["flag_et"],"🇪🇺":["flag_eu"],"🇫🇮":["flag_fi"],"🇫🇯":["flag_fj"],"🇫🇰":["flag_fk"],"🇫🇲":["flag_fm"],"🇫🇴":["flag_fo"],"🇫🇷":["flag_fr"],"🇬🇦":["flag_ga"],"🇬🇧":["flag_gb"],"🇬🇩":["flag_gd"],"🇬🇪":["flag_ge"],"🇬🇫":["flag_gf"],"🇬🇬":["flag_gg"],"🇬🇭":["flag_gh"],"🇬🇮":["flag_gi"],"🇬🇱":["flag_gl"],"🇬🇲":["flag_gm"],"🇬🇳":["flag_gn"],"🇬🇵":["flag_gp"],"🇬🇶":["flag_gq"],"🇬🇷":["flag_gr"],"🇬🇸":["flag_gs"],"🇬🇹":["flag_gt"],"🇬🇺":["flag_gu"],"🇬🇼":["flag_gw"],"🇬🇾":["flag_gy"],"🇭🇰":["flag_hk"],"🇭🇲":["flag_hm"],"🇭🇳":["flag_hn"],"🇭🇷":["flag_hr"],"🇭🇹":["flag_ht"],"🇭🇺":["flag_hu"],"🇮🇨":["flag_ic"],"🇮🇩":["flag_id"],"🇮🇪":["flag_ie"],"🇮🇱":["flag_il"],"🇮🇲":["flag_im"],"🇮🇳":["flag_in"],"🇮🇴":["flag_io"],"🇮🇶":["flag_iq"],"🇮🇷":["flag_ir"],"🇮🇸":["flag_is"],"🇮🇹":["flag_it"],"🇯🇪":["flag_je"],"🇯🇲":["flag_jm"],"🇯🇴":["flag_jo"],"🇯🇵":["flag_jp"],"🇰🇪":["flag_ke"],"🇰🇬":["flag_kg"],"🇰🇭":["flag_kh"],"🇰🇮":["flag_ki"],"🇰🇲":["flag_km"],"🇰🇳":["flag_kn"],"🇰🇵":["flag_kp"],"🇰🇷":["flag_kr"],"🇰🇼":["flag_kw"],"🇰🇾":["flag_ky"],"🇰🇿":["flag_kz"],"🇱🇦":["flag_la"],"🇱🇧":["flag_lb"],"🇱🇨":["flag_lc"],"🇱🇮":["flag_li"],"🇱🇰":["flag_lk"],"🇱🇷":["flag_lr"],"🇱🇸":["flag_ls"],"🇱🇹":["flag_lt"],"🇱🇺":["flag_lu"],"🇱🇻":["flag_lv"],"🇱🇾":["flag_ly"],"🇲🇦":["flag_ma"],"🇲🇨":["flag_mc"],"🇲🇩":["flag_md"],"🇲🇪":["flag_me"],"🇲🇫":["flag_mf"],"🇲🇬":["flag_mg"],"🇲🇭":["flag_mh"],"🇲🇰":["flag_mk"],"🇲🇱":["flag_ml"],"🇲🇲":["flag_mm"],"🇲🇳":["flag_mn"],"🇲🇴":["flag_mo"],"🇲🇵":["flag_mp"],"🇲🇶":["flag_mq"],"🇲🇷":["flag_mr"],"🇲🇸":["flag_ms"],"🇲🇹":["flag_mt"],"🇲🇺":["flag_mu"],"🇲🇻":["flag_mv"],"🇲🇼":["flag_mw"],"🇲🇽":["flag_mx"],"🇲🇾":["flag_my"],"🇲🇿":["flag_mz"],"🇳🇦":["flag_na"],"🇳🇨":["flag_nc"],"🇳🇪":["flag_ne"],"🇳🇫":["flag_nf"],"🇳🇬":["flag_ng"],"🇳🇮":["flag_ni"],"🇳🇱":["flag_nl"],"🇳🇴":["flag_no"],"🇳🇵":["flag_np"],"🇳🇷":["flag_nr"],"🇳🇺":["flag_nu"],"🇳🇿":["flag_nz"],"🇴🇲":["flag_om"],"🇵🇦":["flag_pa"],"🇵🇪":["flag_pe"],"🇵🇫":["flag_pf"],"🇵🇬":["flag_pg"],"🇵🇭":["flag_ph"],"🇵🇰":["flag_pk"],"🇵🇱":["flag_pl"],"🇵🇲":["flag_pm"],"🇵🇳":["flag_pn"],"🇵🇷":["flag_pr"],"🇵🇸":["flag_ps"],"🇵🇹":["flag_pt"],"🇵🇼":["flag_pw"],"🇵🇾":["flag_py"],"🇶🇦":["flag_qa"],"🇷🇪":["flag_re"],"🇷🇴":["flag_ro"],"🇷🇸":["flag_rs"],"🇷🇺":["flag_ru"],"🇷🇼":["flag_rw"],"🇸🇦":["flag_sa"],"🇸🇧":["flag_sb"],"🇸🇨":["flag_sc"],"🇸🇩":["flag_sd"],"🇸🇪":["flag_se"],"🇸🇬":["flag_sg"],"🇸🇭":["flag_sh"],"🇸🇮":["flag_si"],"🇸🇯":["flag_sj"],"🇸🇰":["flag_sk"],"🇸🇱":["flag_sl"],"🇸🇲":["flag_sm"],"🇸🇳":["flag_sn"],"🇸🇴":["flag_so"],"🇸🇷":["flag_sr"],"🇸🇸":["flag_ss"],"🇸🇹":["flag_st"],"🇸🇻":["flag_sv"],"🇸🇽":["flag_sx"],"🇸🇾":["flag_sy"],"🇸🇿":["flag_sz"],"🇹🇦":["flag_ta"],"🇹🇨":["flag_tc"],"🇹🇩":["flag_td"],"🇹🇫":["flag_tf"],"🇹🇬":["flag_tg"],"🇹🇭":["flag_th"],"🇹🇯":["flag_tj"],"🇹🇰":["flag_tk"],"🇹🇱":["flag_tl"],"🇹🇲":["flag_tm"],"🇹🇳":["flag_tn"],"🇹🇴":["flag_to"],"🇹🇷":["flag_tr"],"🇹🇹":["flag_tt"],"🇹🇻":["flag_tv"],"🇹🇼":["flag_tw"],"🇹🇿":["flag_tz"],"🇺🇦":["flag_ua"],"🇺🇬":["flag_ug"],"🇺🇲":["flag_um"],"🇺🇳":["flag_un"],"🇺🇸":["flag_us","usa"],"🇺🇾":["flag_uy"],"🇺🇿":["flag_uz"],"🇻🇦":["flag_va"],"🇻🇨":["flag_vc"],"🇻🇪":["flag_ve"],"🇻🇬":["flag_vg"],"🇻🇮":["flag_vi"],"🇻🇳":["flag_vn"],"🇻🇺":["flag_vu"],"🇼🇫":["flag_wf"],"🇼🇸":["flag_ws"],"🇽🇰":["flag_xk"],"🇾🇪":["flag_ye"],"🇾🇹":["flag_yt"],"🇿🇦":["flag_za"],"🇿🇲":["flag_zm"],"🇿🇼":["flag_zw"],"🏴󠁧󠁢󠁥󠁮󠁧󠁿":["flag_gbeng","england"],"🏴󠁧󠁢󠁳󠁣󠁴󠁿":["flag_gbsct","scotland"],"🏴󠁧󠁢󠁷󠁬󠁳󠁿":["flag_gbwls","wales"]}')
}, function (e, d) {
  e.exports = flarum.core.compat.app
}, function (e, d, a) {
  ! function () {
    var d = ["direction", "boxSizing", "width", "height", "overflowX", "overflowY", "borderTopWidth", "borderRightWidth", "borderBottomWidth", "borderLeftWidth", "borderStyle", "paddingTop", "paddingRight", "paddingBottom", "paddingLeft", "fontStyle", "fontVariant", "fontWeight", "fontStretch", "fontSize", "fontSizeAdjust", "lineHeight", "fontFamily", "textAlign", "textTransform", "textIndent", "textDecoration", "letterSpacing", "wordSpacing", "tabSize", "MozTabSize"],
      a = "undefined" != typeof window,
      u = a && null != window.mozInnerScreenX;

    function o(e, o, n) {
      if (!a) throw new Error("textarea-caret-position#getCaretCoordinates should only be called in a browser");
      var t = n && n.debug || !1;
      if (t) {
        var r = document.querySelector("#input-textarea-caret-position-mirror-div");
        r && r.parentNode.removeChild(r)
      }
      var i = document.createElement("div");
      i.id = "input-textarea-caret-position-mirror-div", document.body.appendChild(i);
      var c = i.style,
        f = window.getComputedStyle ? window.getComputedStyle(e) : e.currentStyle,
        l = "INPUT" === e.nodeName;
      c.whiteSpace = "pre-wrap", l || (c.wordWrap = "break-word"), c.position = "absolute", t || (c.visibility = "hidden"), d.forEach((function (e) {
        l && "lineHeight" === e ? c.lineHeight = f.height : c[e] = f[e]
      })), u ? e.scrollHeight > parseInt(f.height) && (c.overflowY = "scroll") : c.overflow = "hidden", i.textContent = e.value.substring(0, o), l && (i.textContent = i.textContent.replace(/\s/g, " "));
      var _ = document.createElement("span");
      _.textContent = e.value.substring(o) || ".", i.appendChild(_);
      var s = {
        top: _.offsetTop + parseInt(f.borderTopWidth),
        left: _.offsetLeft + parseInt(f.borderLeftWidth),
        height: parseInt(f.lineHeight)
      };
      return t ? _.style.backgroundColor = "#aaa" : document.body.removeChild(i), s
    }
    void 0 !== e.exports ? e.exports = o : a && (window.getCaretCoordinates = o)
  }()
}, function (e, d) {
  e.exports = flarum.core.compat["components/ComposerBody"]
}, function (e, d) {
  e.exports = flarum.core.compat["components/TextEditor"]
}, function (e, d) {
  e.exports = flarum.core.compat["components/TextEditorButton"]
}, function (e, d) {
  e.exports = flarum.core.compat["utils/KeyboardNavigatable"]
}, function (e, d) {
  e.exports = flarum.core.compat.Component
}, function (e, d) {
  e.exports = flarum.core.compat["models/Post"]
}, function (e, d, a) {
  "use strict";
  a.r(d);
  var u = a(2),
    o = a.n(u),
    n = a(3),
    t = a.n(n),
    r = a(1),
    i = a(0),
    c = a(4),
    f = a.n(c),
    l = a(5),
    _ = a.n(l),
    s = a(6),
    g = a.n(s),
    p = function () {
      var e = {
          base: "https://twemoji.maxcdn.com/v/12.1.2/",
          ext: ".png",
          size: "72x72",
          className: "emoji",
          convert: {
            fromCodePoint: function (e) {
              var d = "string" == typeof e ? parseInt(e, 16) : e;
              if (d < 65536) return r(d);
              return r(55296 + ((d -= 65536) >> 10), 56320 + (1023 & d))
            },
            toCodePoint: p
          },
          onerror: function () {
            this.parentNode && this.parentNode.replaceChild(i(this.alt, !1), this)
          },
          parse: function (d, a) {
            a && "function" != typeof a || (a = {
              callback: a
            });
            return ("string" == typeof d ? _ : l)(d, {
              callback: a.callback || c,
              attributes: "function" == typeof a.attributes ? a.attributes : g,
              base: "string" == typeof a.base ? a.base : e.base,
              ext: a.ext || e.ext,
              size: a.folder || (u = a.size || e.size, "number" == typeof u ? u + "x" + u : u),
              className: a.className || e.className,
              onerror: a.onerror || e.onerror
            });
            var u
          },
          replace: m,
          test: function (e) {
            a.lastIndex = 0;
            var d = a.test(e);
            return a.lastIndex = 0, d
          }
        },
        d = {
          "&": "&amp;",
          "<": "&lt;",
          ">": "&gt;",
          "'": "&#39;",
          '"': "&quot;"
        },
        a = /(?:\ud83d\udc68\ud83c\udffc\u200d\ud83e\udd1d\u200d\ud83d\udc68\ud83c\udffb|\ud83d\udc68\ud83c\udffd\u200d\ud83e\udd1d\u200d\ud83d\udc68\ud83c[\udffb\udffc]|\ud83d\udc68\ud83c\udffe\u200d\ud83e\udd1d\u200d\ud83d\udc68\ud83c[\udffb-\udffd]|\ud83d\udc68\ud83c\udfff\u200d\ud83e\udd1d\u200d\ud83d\udc68\ud83c[\udffb-\udffe]|\ud83d\udc69\ud83c\udffb\u200d\ud83e\udd1d\u200d\ud83d\udc68\ud83c[\udffc-\udfff]|\ud83d\udc69\ud83c\udffc\u200d\ud83e\udd1d\u200d\ud83d\udc68\ud83c[\udffb\udffd-\udfff]|\ud83d\udc69\ud83c\udffc\u200d\ud83e\udd1d\u200d\ud83d\udc69\ud83c\udffb|\ud83d\udc69\ud83c\udffd\u200d\ud83e\udd1d\u200d\ud83d\udc68\ud83c[\udffb\udffc\udffe\udfff]|\ud83d\udc69\ud83c\udffd\u200d\ud83e\udd1d\u200d\ud83d\udc69\ud83c[\udffb\udffc]|\ud83d\udc69\ud83c\udffe\u200d\ud83e\udd1d\u200d\ud83d\udc68\ud83c[\udffb-\udffd\udfff]|\ud83d\udc69\ud83c\udffe\u200d\ud83e\udd1d\u200d\ud83d\udc69\ud83c[\udffb-\udffd]|\ud83d\udc69\ud83c\udfff\u200d\ud83e\udd1d\u200d\ud83d\udc68\ud83c[\udffb-\udffe]|\ud83d\udc69\ud83c\udfff\u200d\ud83e\udd1d\u200d\ud83d\udc69\ud83c[\udffb-\udffe]|\ud83e\uddd1\ud83c\udffb\u200d\ud83e\udd1d\u200d\ud83e\uddd1\ud83c\udffb|\ud83e\uddd1\ud83c\udffc\u200d\ud83e\udd1d\u200d\ud83e\uddd1\ud83c[\udffb\udffc]|\ud83e\uddd1\ud83c\udffd\u200d\ud83e\udd1d\u200d\ud83e\uddd1\ud83c[\udffb-\udffd]|\ud83e\uddd1\ud83c\udffe\u200d\ud83e\udd1d\u200d\ud83e\uddd1\ud83c[\udffb-\udffe]|\ud83e\uddd1\ud83c\udfff\u200d\ud83e\udd1d\u200d\ud83e\uddd1\ud83c[\udffb-\udfff]|\ud83e\uddd1\u200d\ud83e\udd1d\u200d\ud83e\uddd1|\ud83d\udc6b\ud83c[\udffb-\udfff]|\ud83d\udc6c\ud83c[\udffb-\udfff]|\ud83d\udc6d\ud83c[\udffb-\udfff]|\ud83d[\udc6b-\udc6d])|(?:\ud83d[\udc68\udc69])(?:\ud83c[\udffb-\udfff])?\u200d(?:\u2695\ufe0f|\u2696\ufe0f|\u2708\ufe0f|\ud83c[\udf3e\udf73\udf93\udfa4\udfa8\udfeb\udfed]|\ud83d[\udcbb\udcbc\udd27\udd2c\ude80\ude92]|\ud83e[\uddaf-\uddb3\uddbc\uddbd])|(?:\ud83c[\udfcb\udfcc]|\ud83d[\udd74\udd75]|\u26f9)((?:\ud83c[\udffb-\udfff]|\ufe0f)\u200d[\u2640\u2642]\ufe0f)|(?:\ud83c[\udfc3\udfc4\udfca]|\ud83d[\udc6e\udc71\udc73\udc77\udc81\udc82\udc86\udc87\ude45-\ude47\ude4b\ude4d\ude4e\udea3\udeb4-\udeb6]|\ud83e[\udd26\udd35\udd37-\udd39\udd3d\udd3e\uddb8\uddb9\uddcd-\uddcf\uddd6-\udddd])(?:\ud83c[\udffb-\udfff])?\u200d[\u2640\u2642]\ufe0f|(?:\ud83d\udc68\u200d\u2764\ufe0f\u200d\ud83d\udc8b\u200d\ud83d\udc68|\ud83d\udc68\u200d\ud83d\udc68\u200d\ud83d\udc66\u200d\ud83d\udc66|\ud83d\udc68\u200d\ud83d\udc68\u200d\ud83d\udc67\u200d\ud83d[\udc66\udc67]|\ud83d\udc68\u200d\ud83d\udc69\u200d\ud83d\udc66\u200d\ud83d\udc66|\ud83d\udc68\u200d\ud83d\udc69\u200d\ud83d\udc67\u200d\ud83d[\udc66\udc67]|\ud83d\udc69\u200d\u2764\ufe0f\u200d\ud83d\udc8b\u200d\ud83d[\udc68\udc69]|\ud83d\udc69\u200d\ud83d\udc69\u200d\ud83d\udc66\u200d\ud83d\udc66|\ud83d\udc69\u200d\ud83d\udc69\u200d\ud83d\udc67\u200d\ud83d[\udc66\udc67]|\ud83d\udc68\u200d\u2764\ufe0f\u200d\ud83d\udc68|\ud83d\udc68\u200d\ud83d\udc66\u200d\ud83d\udc66|\ud83d\udc68\u200d\ud83d\udc67\u200d\ud83d[\udc66\udc67]|\ud83d\udc68\u200d\ud83d\udc68\u200d\ud83d[\udc66\udc67]|\ud83d\udc68\u200d\ud83d\udc69\u200d\ud83d[\udc66\udc67]|\ud83d\udc69\u200d\u2764\ufe0f\u200d\ud83d[\udc68\udc69]|\ud83d\udc69\u200d\ud83d\udc66\u200d\ud83d\udc66|\ud83d\udc69\u200d\ud83d\udc67\u200d\ud83d[\udc66\udc67]|\ud83d\udc69\u200d\ud83d\udc69\u200d\ud83d[\udc66\udc67]|\ud83c\udff3\ufe0f\u200d\u26a7\ufe0f|\ud83c\udff3\ufe0f\u200d\ud83c\udf08|\ud83c\udff4\u200d\u2620\ufe0f|\ud83d\udc15\u200d\ud83e\uddba|\ud83d\udc41\u200d\ud83d\udde8|\ud83d\udc68\u200d\ud83d[\udc66\udc67]|\ud83d\udc69\u200d\ud83d[\udc66\udc67]|\ud83d\udc6f\u200d\u2640\ufe0f|\ud83d\udc6f\u200d\u2642\ufe0f|\ud83e\udd3c\u200d\u2640\ufe0f|\ud83e\udd3c\u200d\u2642\ufe0f|\ud83e\uddde\u200d\u2640\ufe0f|\ud83e\uddde\u200d\u2642\ufe0f|\ud83e\udddf\u200d\u2640\ufe0f|\ud83e\udddf\u200d\u2642\ufe0f)|[#*0-9]\ufe0f?\u20e3|(?:[©®\u2122\u265f]\ufe0f)|(?:\ud83c[\udc04\udd70\udd71\udd7e\udd7f\ude02\ude1a\ude2f\ude37\udf21\udf24-\udf2c\udf36\udf7d\udf96\udf97\udf99-\udf9b\udf9e\udf9f\udfcd\udfce\udfd4-\udfdf\udff3\udff5\udff7]|\ud83d[\udc3f\udc41\udcfd\udd49\udd4a\udd6f\udd70\udd73\udd76-\udd79\udd87\udd8a-\udd8d\udda5\udda8\uddb1\uddb2\uddbc\uddc2-\uddc4\uddd1-\uddd3\udddc-\uddde\udde1\udde3\udde8\uddef\uddf3\uddfa\udecb\udecd-\udecf\udee0-\udee5\udee9\udef0\udef3]|[\u203c\u2049\u2139\u2194-\u2199\u21a9\u21aa\u231a\u231b\u2328\u23cf\u23ed-\u23ef\u23f1\u23f2\u23f8-\u23fa\u24c2\u25aa\u25ab\u25b6\u25c0\u25fb-\u25fe\u2600-\u2604\u260e\u2611\u2614\u2615\u2618\u2620\u2622\u2623\u2626\u262a\u262e\u262f\u2638-\u263a\u2640\u2642\u2648-\u2653\u2660\u2663\u2665\u2666\u2668\u267b\u267f\u2692-\u2697\u2699\u269b\u269c\u26a0\u26a1\u26a7\u26aa\u26ab\u26b0\u26b1\u26bd\u26be\u26c4\u26c5\u26c8\u26cf\u26d1\u26d3\u26d4\u26e9\u26ea\u26f0-\u26f5\u26f8\u26fa\u26fd\u2702\u2708\u2709\u270f\u2712\u2714\u2716\u271d\u2721\u2733\u2734\u2744\u2747\u2757\u2763\u2764\u27a1\u2934\u2935\u2b05-\u2b07\u2b1b\u2b1c\u2b50\u2b55\u3030\u303d\u3297\u3299])(?:\ufe0f|(?!\ufe0e))|(?:(?:\ud83c[\udfcb\udfcc]|\ud83d[\udd74\udd75\udd90]|[\u261d\u26f7\u26f9\u270c\u270d])(?:\ufe0f|(?!\ufe0e))|(?:\ud83c[\udf85\udfc2-\udfc4\udfc7\udfca]|\ud83d[\udc42\udc43\udc46-\udc50\udc66-\udc69\udc6e\udc70-\udc78\udc7c\udc81-\udc83\udc85-\udc87\udcaa\udd7a\udd95\udd96\ude45-\ude47\ude4b-\ude4f\udea3\udeb4-\udeb6\udec0\udecc]|\ud83e[\udd0f\udd18-\udd1c\udd1e\udd1f\udd26\udd30-\udd39\udd3d\udd3e\uddb5\uddb6\uddb8\uddb9\uddbb\uddcd-\uddcf\uddd1-\udddd]|[\u270a\u270b]))(?:\ud83c[\udffb-\udfff])?|(?:\ud83c\udff4\udb40\udc67\udb40\udc62\udb40\udc65\udb40\udc6e\udb40\udc67\udb40\udc7f|\ud83c\udff4\udb40\udc67\udb40\udc62\udb40\udc73\udb40\udc63\udb40\udc74\udb40\udc7f|\ud83c\udff4\udb40\udc67\udb40\udc62\udb40\udc77\udb40\udc6c\udb40\udc73\udb40\udc7f|\ud83c\udde6\ud83c[\udde8-\uddec\uddee\uddf1\uddf2\uddf4\uddf6-\uddfa\uddfc\uddfd\uddff]|\ud83c\udde7\ud83c[\udde6\udde7\udde9-\uddef\uddf1-\uddf4\uddf6-\uddf9\uddfb\uddfc\uddfe\uddff]|\ud83c\udde8\ud83c[\udde6\udde8\udde9\uddeb-\uddee\uddf0-\uddf5\uddf7\uddfa-\uddff]|\ud83c\udde9\ud83c[\uddea\uddec\uddef\uddf0\uddf2\uddf4\uddff]|\ud83c\uddea\ud83c[\udde6\udde8\uddea\uddec\udded\uddf7-\uddfa]|\ud83c\uddeb\ud83c[\uddee-\uddf0\uddf2\uddf4\uddf7]|\ud83c\uddec\ud83c[\udde6\udde7\udde9-\uddee\uddf1-\uddf3\uddf5-\uddfa\uddfc\uddfe]|\ud83c\udded\ud83c[\uddf0\uddf2\uddf3\uddf7\uddf9\uddfa]|\ud83c\uddee\ud83c[\udde8-\uddea\uddf1-\uddf4\uddf6-\uddf9]|\ud83c\uddef\ud83c[\uddea\uddf2\uddf4\uddf5]|\ud83c\uddf0\ud83c[\uddea\uddec-\uddee\uddf2\uddf3\uddf5\uddf7\uddfc\uddfe\uddff]|\ud83c\uddf1\ud83c[\udde6-\udde8\uddee\uddf0\uddf7-\uddfb\uddfe]|\ud83c\uddf2\ud83c[\udde6\udde8-\udded\uddf0-\uddff]|\ud83c\uddf3\ud83c[\udde6\udde8\uddea-\uddec\uddee\uddf1\uddf4\uddf5\uddf7\uddfa\uddff]|\ud83c\uddf4\ud83c\uddf2|\ud83c\uddf5\ud83c[\udde6\uddea-\udded\uddf0-\uddf3\uddf7-\uddf9\uddfc\uddfe]|\ud83c\uddf6\ud83c\udde6|\ud83c\uddf7\ud83c[\uddea\uddf4\uddf8\uddfa\uddfc]|\ud83c\uddf8\ud83c[\udde6-\uddea\uddec-\uddf4\uddf7-\uddf9\uddfb\uddfd-\uddff]|\ud83c\uddf9\ud83c[\udde6\udde8\udde9\uddeb-\udded\uddef-\uddf4\uddf7\uddf9\uddfb\uddfc\uddff]|\ud83c\uddfa\ud83c[\udde6\uddec\uddf2\uddf3\uddf8\uddfe\uddff]|\ud83c\uddfb\ud83c[\udde6\udde8\uddea\uddec\uddee\uddf3\uddfa]|\ud83c\uddfc\ud83c[\uddeb\uddf8]|\ud83c\uddfd\ud83c\uddf0|\ud83c\uddfe\ud83c[\uddea\uddf9]|\ud83c\uddff\ud83c[\udde6\uddf2\uddfc]|\ud83c[\udccf\udd8e\udd91-\udd9a\udde6-\uddff\ude01\ude32-\ude36\ude38-\ude3a\ude50\ude51\udf00-\udf20\udf2d-\udf35\udf37-\udf7c\udf7e-\udf84\udf86-\udf93\udfa0-\udfc1\udfc5\udfc6\udfc8\udfc9\udfcf-\udfd3\udfe0-\udff0\udff4\udff8-\udfff]|\ud83d[\udc00-\udc3e\udc40\udc44\udc45\udc51-\udc65\udc6a-\udc6d\udc6f\udc79-\udc7b\udc7d-\udc80\udc84\udc88-\udca9\udcab-\udcfc\udcff-\udd3d\udd4b-\udd4e\udd50-\udd67\udda4\uddfb-\ude44\ude48-\ude4a\ude80-\udea2\udea4-\udeb3\udeb7-\udebf\udec1-\udec5\uded0-\uded2\uded5\udeeb\udeec\udef4-\udefa\udfe0-\udfeb]|\ud83e[\udd0d\udd0e\udd10-\udd17\udd1d\udd20-\udd25\udd27-\udd2f\udd3a\udd3c\udd3f-\udd45\udd47-\udd71\udd73-\udd76\udd7a-\udda2\udda5-\uddaa\uddae-\uddb4\uddb7\uddba\uddbc-\uddca\uddd0\uddde-\uddff\ude70-\ude73\ude78-\ude7a\ude80-\ude82\ude90-\ude95]|[\u23e9-\u23ec\u23f0\u23f3\u267e\u26ce\u2705\u2728\u274c\u274e\u2753-\u2755\u2795-\u2797\u27b0\u27bf\ue50a])|\ufe0f/g,
        u = /\uFE0F/g,
        o = String.fromCharCode(8205),
        n = /[&<>'"]/g,
        t = /^(?:iframe|noframes|noscript|script|select|style|textarea)$/,
        r = String.fromCharCode;
      return e;

      function i(e, d) {
        return document.createTextNode(d ? e.replace(u, "") : e)
      }

      function c(e, d) {
        return "".concat(d.base, d.size, "/", e, d.ext)
      }

      function f(e) {
        return p(e.indexOf(o) < 0 ? e.replace(u, "") : e)
      }

      function l(e, d) {
        for (var u, o, n, r, c, l, _, s, g, m, p, b, h, w = function e(d, a) {
            for (var u, o, n = d.childNodes, r = n.length; r--;) 3 === (o = (u = n[r]).nodeType) ? a.push(u) : 1 !== o || "ownerSVGElement" in u || t.test(u.nodeName.toLowerCase()) || e(u, a);
            return a
          }(e, []), k = w.length; k--;) {
          for (n = !1, r = document.createDocumentFragment(), l = (c = w[k]).nodeValue, s = 0; _ = a.exec(l);) {
            if ((g = _.index) !== s && r.appendChild(i(l.slice(s, g), !0)), b = f(p = _[0]), s = g + p.length, h = d.callback(b, d), b && h) {
              for (o in (m = new Image).onerror = d.onerror, m.setAttribute("draggable", "false"), u = d.attributes(p, b)) u.hasOwnProperty(o) && 0 !== o.indexOf("on") && !m.hasAttribute(o) && m.setAttribute(o, u[o]);
              m.className = d.className, m.alt = p, m.src = h, n = !0, r.appendChild(m)
            }
            m || r.appendChild(i(p, !1)), m = null
          }
          n && (s < l.length && r.appendChild(i(l.slice(s), !0)), c.parentNode.replaceChild(r, c))
        }
        return e
      }

      function _(e, d) {
        return m(e, (function (e) {
          var a, u, o = e,
            t = f(e),
            r = d.callback(t, d);
          if (t && r) {
            for (u in o = "<img ".concat('class="', d.className, '" ', 'draggable="false" ', 'alt="', e, '"', ' src="', r, '"'), a = d.attributes(e, t)) a.hasOwnProperty(u) && 0 !== u.indexOf("on") && -1 === o.indexOf(" " + u + "=") && (o = o.concat(" ", u, '="', a[u].replace(n, s), '"'));
            o = o.concat("/>")
          }
          return o
        }))
      }

      function s(e) {
        return d[e]
      }

      function g() {
        return null
      }

      function m(e, d) {
        return String(e).replace(a, d)
      }

      function p(e, d) {
        for (var a = [], u = 0, o = 0, n = 0; n < e.length;) u = e.charCodeAt(n++), o ? (a.push((65536 + (o - 55296 << 10) + (u - 56320)).toString(16)), o = 0) : 55296 <= u && u <= 56319 ? o = u : a.push(u.toString(16));
        return a.join(d || "-")
      }
    }(),
    b = String.fromCharCode(8205),
    h = /\uFE0F/g;

  function w(e) {
    return p.convert.toCodePoint(e.indexOf(b) < 0 ? e.replace(h, "") : e)
  }
  var k = a(7),
    y = a.n(k);
  var v = a(8),
    x = function (e) {
      var d, a;

      function u() {
        return e.apply(this, arguments) || this
      }
      a = e, (d = u).prototype = Object.create(a.prototype), d.prototype.constructor = d, d.__proto__ = a;
      var o = u.prototype;
      return o.init = function () {
        this.active = !1, this.index = 0, this.keyWasJustPressed = !1
      }, o.view = function () {
        return m("ul", {
          className: "Dropdown-menu EmojiDropdown"
        }, m("li", {
          className: "Dropdown-header",
          key: "0"
        }, app.translator.trans("flarum-emoji.forum.composer.type_to_search_text")), this.props.items.map((function (e) {
          return m("li", {
            key: e.attrs.key
          }, e)
        })))
      }, o.show = function (e, d) {
        this.$().show().css({
          left: e + "px",
          top: d + "px"
        }), this.active = !0
      }, o.hide = function () {
        this.$().hide(), this.active = !1
      }, o.navigate = function (e) {
        var d = this;
        this.keyWasJustPressed = !0, this.setIndex(this.index + e, !0), clearTimeout(this.keyWasJustPressedTimeout), this.keyWasJustPressedTimeout = setTimeout((function () {
          return d.keyWasJustPressed = !1
        }), 500)
      }, o.complete = function () {
        this.$("li:not(.Dropdown-header)").eq(this.index).find("button").click()
      }, o.setIndex = function (e, d) {
        if (!this.keyWasJustPressed || d) {
          var a = this.$(),
            u = a.find("li:not(.Dropdown-header)"),
            o = e;
          o < 0 ? o = u.length - 1 : o >= u.length && (o = 0), this.index = o;
          var n = u.removeClass("active").eq(o).addClass("active");
          if (d) {
            var t, r = a.scrollTop(),
              i = a.offset().top,
              c = i + a.outerHeight(),
              f = n.offset().top,
              l = f + n.outerHeight();
            f < i ? t = r - i + f - parseInt(a.css("padding-top"), 10) : l > c && (t = r - c + l + parseInt(a.css("padding-bottom"), 10)), void 0 !== t && a.stop(!0).animate({
              scrollTop: t
            }, 100)
          }
        }
      }, u
    }(a.n(v).a);

  function j() {
    var e = Object.keys(r);
    Object(i.extend)(f.a.prototype, "config", (function (d, a) {
      if (!a) {
        var u, o, n = this,
          i = $('<div class="ComposerBody-emojiDropdownContainer"></div>'),
          c = new x({
            items: []
          }),
          f = this.$("textarea").wrap('<div class="ComposerBody-emojiWrapper"></div>');
        this.navigator = new y.a, this.navigator.when((function () {
          return c.active
        })).onUp((function () {
          return c.navigate(-1)
        })).onDown((function () {
          return c.navigate(1)
        })).onSelect(c.complete.bind(c)).onCancel(c.hide.bind(c)).bindTo(f), f.after(i).on("click keyup input", (function (d) {
          var a = this;
          if (-1 === [9, 13, 27, 40, 38, 37, 39].indexOf(d.which)) {
            var l = this.selectionStart;
            if (!(this.selectionEnd - l > 0)) {
              var _ = this.value;
              u = 0;
              for (var s = l - 1; s >= 0; s--) {
                var g = _.substr(s, 1);
                if (!/[a-z0-9]|\+|\-|_|\:/.test(g)) break;
                if (":" === g && (0 == s || /\s/.test(_.substr(s - 1, 1)))) {
                  u = s + 1;
                  break
                }
              }
              if (c.hide(), c.active = !1, u) {
                o = _.substring(u, l).toLowerCase();
                var p = function (e) {
                  var d = e.emoji,
                    a = e.name,
                    o = e.code;
                  return m("button", {
                    key: d,
                    onclick: function () {
                      return a = u - 1 + (e = d + " ").length, o = n.content(), n.editor.setValue(o.substring(0, u - 1) + e + o.substr(f[0].selectionStart)), n.editor.setSelectionRange(a, a), void c.hide();
                      var e, a, o
                    },
                    onmouseenter: function () {
                      c.setIndex($(this).parent().index() - 1)
                    }
                  }, m("img", {
                    alt: d,
                    class: "emoji",
                    draggable: "false",
                    src: "//cdn.jsdelivr.net/gh/twitter/twemoji@12/assets/72x72/" + o + ".png"
                  }), a)
                };
                ! function () {
                  var d, n, f = [],
                    l = (d = o, n = new RegExp("\\(([" + "+.*?[]{}()^$|\\".replace(/(.)/g, "\\$1") + "])\\)", "g"), new RegExp("(.*)" + d.toLowerCase().replace(/(.)/g, "($1)(.*?)").replace(n, "(\\$1)") + "$", "i")),
                    _ = 7,
                    s = function (d) {
                      for (var a = 0; a < e.length && _ > 0; a++) {
                        var u = e[a];
                        if (-1 === f.indexOf(u)) {
                          var o = r[u],
                            n = Array.isArray(o),
                            t = 0;
                          for (o = n ? o : o[Symbol.iterator]();;) {
                            var i;
                            if (n) {
                              if (t >= o.length) break;
                              i = o[t++]
                            } else {
                              if ((t = o.next()).done) break;
                              i = t.value
                            }
                            if (d(i)) {
                              --_, f.push(u);
                              break
                            }
                          }
                        }
                      }
                    };
                  s((function (e) {
                    return 0 === e.indexOf(o)
                  })), s((function (e) {
                    return l.test(e)
                  }));
                  var g = f.map((function (e) {
                    return {
                      emoji: e,
                      name: r[e][0],
                      code: w(e)
                    }
                  })).map(p);
                  if (g.length) {
                    c.props.items = g, m.render(i[0], c.render()), c.show();
                    var b = t()(a, u),
                      h = c.$().outerWidth(),
                      k = c.$().outerHeight(),
                      y = c.$().offsetParent(),
                      v = b.left,
                      x = b.top + 15;
                    x + k > y.height() && (x = b.top - k - 15), v + h > y.width() && (v = y.width() - h), x = Math.max(-$(a).offset().top, x), v = Math.max(-$(a).offset().left, v), c.show(v, x)
                  }
                }(), c.setIndex(0), c.$().scrollTop(0), c.active = !0
              }
            }
          }
        }))
      }
    })), Object(i.extend)(_.a.prototype, "toolbarItems", (function (e) {
      var d = this;
      e.add("emoji", m(g.a, {
        onclick: function () {
          return d.insertAtCursor(":")
        },
        icon: "far fa-smile"
      }, app.translator.trans("flarum-emoji.forum.composer.emoji_tooltip")))
    }))
  }
  var z = a(9),
    C = a.n(z);
  o.a.initializers.add("flarum-emoji", (function () {
    j(), Object(i.override)(C.a.prototype, "contentHtml", (function (e) {
      var d = e();
      return this.oldContentHtml !== d && (this.emojifiedContentHtml = p.parse(d), this.oldContentHtml = d), this.emojifiedContentHtml
    })), Object(i.override)(s9e.TextFormatter, "preview", (function (e, d, a) {
      e(d, a), p.parse(a)
    }))
  }))
}]);;