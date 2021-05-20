// initial sidebar
window.onload = () => {
  initSidebar();
  window.scroll(0, 0);
  if (location.hash.length > 0) {
    setTimeout(() => {
      scrollTocNav(location.hash.replace('#', ''));
    }, 100)
  }
}

// click event
window.onclick = (event) => {
  let elem = getElemByEvent(event);
  if (elem.id === 'current') {
    toggleTocNav();
  } else {
    if (elem.id === 'list-button') {
      toggleSideAll();
    } else if (elem.id.startsWith('cat_')) {
      toggleSidebar(elem);
    } else if (elem.href) {
      let id = elem.href.split('/').slice(-1)[0];
      if (id.startsWith('#')) {
        setTimeout(() => {
          scrollTocNav(id.replace('#', ''));
        }, 200);
      }
    }
    toggleTocNav(true);
  }
}

// scroll event
window.onscroll = () => {
  const heads = document.querySelectorAll('.markdown h2, .markdown h3');
  const headerY = 64;
  let preId = '';
  let preDiff = 100000;
  let nextId = '';
  let nextDiff = 100000;
  for (let head of heads) {
    const pos = head.getBoundingClientRect().top - headerY;
    if (pos < 0 && Math.abs(0 - pos) < preDiff) {
      preDiff = Math.abs(0 - pos);
      preId = head.id;
    } else if (pos > 0 && Math.abs(0 - pos) < nextDiff) {
      nextDiff = Math.abs(0 - pos);
      nextId = head.id;
    }
  }

  if (!preId) {
    scrollTocNav(nextId);
  } else {
    (preDiff < nextDiff) ? scrollTocNav(preId) : scrollTocNav(nextId);
  }
}

(function() {

// hamburger menu
const $wrapper = document.getElementById('menu');
const $navBtn = document.getElementById('nav-btn');
const $ancorLink = document.querySelectorAll('a[href^="#"]');
$ancorLink.forEach(function(
    button) { button.addEventListener('click', navClose); });

$navBtn.addEventListener('click', navToggle);

function navToggle() {
  if ($wrapper.classList.contains('header__list--open')) {
    navClose();
  } else {
    navOpen();
  }
}

function navOpen() { $wrapper.classList.add('header__list--open'); }

function navClose() { $wrapper.classList.remove('header__list--open'); }

// toc toggle
const tocWrap = document.getElementById('current');
document.addEventListener('click', tocWrap, currentToggle);

function currentToggle() {
  if (tocWrap.classList.contains('open')) {
    tocClose();
  } else {
    tocOpen();
  }
}

function tocOpen() { tocWrap.classList.add('open'); }

function tocClose() { tocWrap.classList.remove('open'); }

// smooth scroll
const headerHight = document.getElementById('header').offsetHeight;

let smoothScroll = (target, offset) => {
  let toY;
  let nowY = window.pageYOffset;
  const divisor = 8;
  const range = (divisor / 2) + 1;

  const targetRect = target.getBoundingClientRect();
  const targetY = targetRect.top + nowY - offset;

  (function() {
    let thisFunc = arguments.callee;
    toY = nowY + Math.round((targetY - nowY) / divisor);
    window.scrollTo(0, toY);
    nowY = toY;

    if (document.body.clientHeight - window.innerHeight < toY) {
      window.scrollTo(0, document.body.clientHeight);
      return;
    }
    if (toY >= targetY + range || toY <= targetY - range) {
      window.setTimeout(thisFunc, 10);

    } else {
      window.scrollTo(0, targetY);
    }
  })();
};

const smoothOffset = headerHight;
const links = document.querySelectorAll('a[href*="#"]');
for (let i = 0; i < links.length; i++) {
  links[i].addEventListener('click', function(e) {
    const href = e.currentTarget.getAttribute('href');
    const splitHref = href.split('#');
    const targetID = splitHref[1];
    const target = document.getElementById(targetID);

    if (target) {
      smoothScroll(target, smoothOffset);
    } else {
      return true;
    }
    return false;
  });
}
})();

// initialize sidebar style
const initSidebar = () => {
  const sidebar = document.getElementById('list-body');
  const paths = window.location.href.split("/").filter((v) => {
      if (v.length != 0) {
          return v;
      }
  });
  if (sidebar) {
    const lastPath = paths[paths.length - 1];
    for (let child of sidebar.children) {
      if ((lastPath == "docs" || lastPath.match("v\[0-9\]+")) && child.className == "withchild") {
        child.className = "withchild open";
      }
      let isOpen = false;
      let category = document.getElementById(child.id);
      const contents = category.getElementsByTagName('li');
      for (const link of contents) {
        if (link.className === 'view') isOpen = !isOpen;
      }
      if (isOpen) {
        category.className = "withchild open"
      }
    }
  }
}

// toggle all by click
const toggleSideAll = () => {
  let sidebar = document.getElementById('list-body');
  let rootBar = document.getElementById('list-button');
  if (sidebar) {
    if (sidebar.style.display.length > 0) {
      sidebar.style.display = '';
      rootBar.className = 'index open';
    } else if (sidebar.style.length === 0) {
      sidebar.style.display = 'none';
      rootBar.className = 'index';
    }
  }
}

// toggle each category by click
const toggleSidebar = (elem) => {
  if (elem.className.includes('open')) {
    elem.className = 'withchild';
  } else {
    elem.className = 'withchild open';
  }
}

// toggle toc nav
const toggleTocNav = (close = false) => {
  let elem = document.getElementById('current');
  if (!elem) return;
  if (close) {
    elem.className = 'current';
  } else {
    if (elem.className.includes('open')) {
      elem.className = 'current'
    } else {
      elem.className = 'current open';
    }
  }
}

// scroll toc nav
const scrollTocNav = (id) => {
  let toc = document.querySelectorAll('.current a');
  id = '#' + id;
  for (const link of toc) {
    link.className = (link.hash === id) ? 'view' : '';
  }
}

const getElemByEvent = (event) => {
  let elem = event.target;
  // for TOC
  if (!elem.id && (elem.className === 'dot' || elem.className === 'menu')) {
    elem = getParentByElem(elem);
    if (!elem.id && elem.className === 'menu') {
      elem = getParentByElem(elem)
    }
  }
  return elem;
}

const getParentByElem = (elem) => {
  return elem.parentNode;
}
