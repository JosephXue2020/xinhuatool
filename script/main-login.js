! function() {
    var e = new Date;
    2020 === e.getFullYear() && 3 === e.getMonth() && (4 !== e.getDate() && 5 !== e.getDate() || $(".loginMainBodyBox").addClass("black-white-theme"))
}(),

$(function() {
    window.self != window.top && window.top.location.replace(GLOBAL.WEBROOT + "/loginxhs"), $(window.top.document.body).hasClass("loginMainBodyBox") || window.top.location.replace(GLOBAL.WEBROOT + "/loginxhs"), require(["base64"], function(e) {
            var i = $("#loginform").attr("action"),
                o = $("#loginAction").val(),
                t = ($("#xinhua_login_service_login_url").val(), $("#xinhua_client_main_url").val()),
                n = ($("#xinhua_client_from").val(), $("#xinhua_login_check_login_url").val()),
                a = "0",
                l = function() {
                    var i = $("#password").val();
                    i.length && ((i.length < 7 || "_pwdtag" != i.substring(i.length - 7, i.length)) && $("#password").val(e.encode($("#password").val() + "_xhs") + "_pwdtag"),
                        $("#loginform").submit())
                },

                r = jstz.determine().name();
            $.ajax({ url: GLOBAL.WEBROOT + "/loginxhs/ssoopen", async: !1, dataType: "json", success: function(e) { a = e } });
            var s = $(".promptBox").find("span"),
                c = s.text();
            "1" == a && "" == c && xinhua_sso.checkLogined("xingonggao", n, t), $.cookie("InvalidOrderFifteenDaysAfter", null), MainLogin.generAdvertise(), window.self != window.top && window.top.location.replace($webroot + "login"), $.validator.addMethod("postCode", function(e, i) { var o = /^\d?$/; return "" == e || !o.test(e) }, "邮编格式不正确"), window.errorInfo = "", $.validator.addMethod("defineRequired", function(e, i) { return errorInfo = $(i).attr("requiredInfo"), "" != e || "" != $.trim(e) }, function(e, i) { return $(i).attr("requiredInfo") });
            GLOBAL.WEBROOT;
            $("#refleshCaptchaCode").click(function() { $("#captchaCodeImg").attr("src", GLOBAL.WEBROOT + "/captcha/CapthcaImage?" + (new Date).getTime()) }), $("#loginsubmit").click(function() { return $("#loginReferer").val(window.location.search), "" == $.trim($("#password").val()) && "" == $.trim($("#userName").val()) ? ($("#loginErrorInfo").empty(), void $("#loginErrorInfo").html("<div class='promptBox'><i class='picon'></i>请输入账户名和密码</div>")) : "" == $.trim($("#userName").val()) ? ($("#loginErrorInfo").empty(), void $("#loginErrorInfo").html("<div class='promptBox'><i class='picon'></i>请输入账号名</div>")) : "" == $.trim($("#password").val()) ? ($("#loginErrorInfo").empty(), void $("#loginErrorInfo").html("<div class='promptBox'><i class='picon'></i>请输入密码</div>")) : ($("#loginform").attr("action", o), $("#timezoneVal").val(r), $.cookie("toolIntro_lastTime", null, { path: "/" }), void l()) }), $("#loginform").keydown(function(e) { 13 == e.keyCode && $("#loginsubmit").click() }),

                $("#visitorlogin").click(function() {
                    var e = $("#anonymous_u").val(),
                        o = $("#anonymous_p").val(),
                        t = $("#anonymous_OpenCode").val();
                    if ($("#loginform").attr("action", i), "1" == t)
                        $("#userName").val(e), $("#password").val(o);
                    else if ("1" == a) {
                        $("#userName").val(e), $("#password").val(o);
                        var n = $("#visitorSSOLogin").val();
                        $("#loginform").attr("action", n), l(), $("#userName").val(""), $("#password").val("")
                    } else
                        $("#userName").val(e), $("#password").val(o), l(), $("#userName").val(""), $("#password").val("")
                });

            $(".menu-active").hover(function(e) { "" == $.trim($(".menu-con-title", $(this)).html()) ? $(this).unbind("hover") : $(".menu-con-content", $(this)).hasClass("mCustomScrollbar") || $(".menu-con-content", $(this)).mCustomScrollbar({ scrollInertia: 150 }), e.stopPropagation(), e.preventDefault() }), pageConfig.config({ plugin: [] })
        }),

        $.ajax({
            type: "POST",
            url: $webroot + "site/currentSite",
            dataType: "json",
            success: function(e) {
                var i = $(".site-item[data-id='" + e + "']");
                i.length && ($("#site-selected").text(i.text()), $(".site-item").removeClass("active"), i.addClass("active"))
            },
            error: function() { console.log("获取站点异常"), $(".site-item[data-id='1']").addClass("active") }
        })
});


var MainLogin = {
    generAdvertise: function() {
        var e = $("#slideBox");
        $.eAjax({
            url: GLOBAL.WEBROOT + "/leaflet/qryLeafletList",
            data: { placeId: "1", placeSize: "5", placeWidth: "1920", placeHeight: "446", status: "1" },
            async: !0,
            type: "post",
            dataType: "json",
            success: function(i, o) {
                if (null != i)
                    if (i && i.respList && i.respList.length > 0) {
                        var t = $(".image-slider-num", e),
                            n = $(".image-slider-img", e);
                        MainLogin.doAdList(t, n, i.respList)
                    } else $(e).empty(), $(e).append("<div class ='pro-empty'></div>")
            },
            error: function() {}
        })
    },
    doAdList: function(e, i, o) {
        i.empty();
        var t = "",
            n = "";
        t += "<ul class='carousel-indicators'>", $.each(o, function(e, i) { t += 0 == e ? "<li class='on'>" : "<li>", t += "<img style='width:100px;height:70px;' src='" + i.vfsUrl + "' alt='" + i.advertiseTitle + "'>", t += "</li>" }), t += "<ul>", e.append(t), i.empty(), n = "<ul>", $.each(o, function(e, i) { n += "<li class='item'>", n += "<img src=" + i.vfsUrl + " alt=''>", n += "</li>" }), n += "</ul>", i.append(n), i = null, o.length >= 2 && $("#slideBox").slide({ mainCell: ".bd ul", autoPlay: !0, effect: "fade", interTime: 6e3 }), e = null;
        var a = $(".login-header").outerHeight() + $(".menu-list").outerHeight() + $(".public-footer").outerHeight(),
            l = $(window).height() - a;
        l = l > 446 ? l : 446, $(".login-panel-banner").height(l), $(".login-panel-banner .bd li").height(l)
    }
};
//# sourceMappingURL=../../../maps/zh_CN/common/login/main-login.js.map