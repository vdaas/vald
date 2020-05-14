$(function() {

  //header
  $('.header-nav__icon').on('click', function() {
    $('.header__list').toggleClass('header__list--open');
    return false;
  });

  $('.header__list li a[href^="#"]').on('click', function() {
    $('.header__list').removeClass('header__list--open');
    return false;
  });


  //contents scroll effecrt
  $(window).scroll(function() {
    $(".mot").each(function() {
      var hh = $(".header").height();
      var position = $(this).offset().top + hh;
      var scroll = $(window).scrollTop();
      var windowHeight = $(window).height();
      if (scroll > position - windowHeight) {
        $(this).addClass('active');
      }
    });
  });

  //scroll top
  $(window).scroll(function() {
    var y = $(window).scrollTop();
    if (y < 500) $(".mdl-totop").slideUp(200);
    else $(".mdl-totop").slideDown(200);
  });
  var pagetop = $('.mdl-totop');
  pagetop.click(function() {
    $('body, html').animate({
      scrollTop: 0
    }, 700);
    return false;
  });

  //#link
  $('a[href^="#"]').click(function() {

    var Hash = $(this.hash);
    var $headerheight = $('.header').height();
    var HashOffset = $(Hash).offset().top;
    var HashOffset_true = HashOffset - $headerheight;

    $("html,body").animate({
      scrollTop: HashOffset_true
    }, 700);

    return false;
  });

});
