// initial sidebar
window.onload = async () => {
  initSidebar();
  const githubObj = await getGitHubStar();
  //console.log(githubObj.stargazers_count);
  if (githubObj.stargazers_count) {
    let elem = document.getElementById("git-star-num");
    elem.innerHTML = githubObj.stargazers_count;
  }

  if (location.hash.length > 0) {
    const height = document.getElementsByClassName('header__nav')[0].offsetHeight;
    window.scrollBy(0, -height);
    setTimeout(() => {
      scrollTocNav(location.hash.replace('#', ''));
    }, 100)
  } else {
    window.scroll(0, 0);
  }
}

// click event
window.onclick = (event) => {
  let elem = getElemByEvent(event);
  // copy code to clipboard
  if (elem.id === '' && elem.className === "highlight") {
    const code = elem.children[0].children[0].innerText;
    navigator.clipboard.writeText(code)
      .then(() => {
        if (!elem.classList.contains("clicked")) {
          elem.classList.add("clicked")
          setTimeout(() => {
            elem.classList.remove("clicked");
          }, 500);
        } else {
          elem.classList.remove("clicked");
        }
      })
      .catch((e) => {
        console.log("failed to copy code");
      })
  } else if (elem.id === 'current') {
    toggleTocNav();
  } else if (elem.className.includes('version__') || elem.className.includes('__version')) {
    if (elem.className == 'header__version' || elem.className == 'version__current') {
      toggleVersion();
    } else {
      setVersion(elem);
    }
  } else {
    if (elem.id === 'list-button') {
      toggleSideAll();
      toggleTocNav(true);
    } else if (elem.id.startsWith('cat_')) {
      toggleSidebar(elem);
      toggleTocNav(true);
    } else if (elem.href) {
      let id = elem.href.split('/').slice(-1)[0];
      if (id.startsWith('#')) {
        setTimeout(() => {
          scrollTocNav(id.replace('#', ''));
        }, 200);
      }
      toggleTocNav(true);
    }
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

function toggleVersion() {
  if (document.getElementById('version_details').open) {
    document.getElementById('version_details').setAttribute('open', true);
  } else {
    document.getElementById('version_details').open = false;
  }
}

// set document version
const setVersion = (elem) => {
  if (elem.text === '' || elem.text === undefined) {
    document.getElementById('version_details').removeAttribute('open');
  } else if (elem.text.startsWith('v')) {
    const beforeVersion = document.getElementById('current_version').textContent.trim();
    document.getElementById('current_version').textContent = elem.text;
    document.getElementById('version_details').removeAttribute('open');
    let url = location.href;
    const nextVersion = elem.className.includes('latest') ? '' : elem.text + '/';
    if (url.includes('/docs/')) {
      let vOfUrl = "";
      if (url.split('/docs/').length > 1) {
        vOfUrl = url.split('/docs/')[1].split('/')[0];
      }
      const regex = /v\d{1}\.\d{1}/;
      const match = vOfUrl.match(regex);
      // move to new document url .
      if (vOfUrl.length > 0) {
        if (vOfUrl === beforeVersion) {
          url = url.replace(beforeVersion + '/', nextVersion);
        } else if (match && match.length === 1) {
          // when 404 page is show, this branch will run.
          url = url.replace(vOfUrl + '/', nextVersion);
        } else {
          url = url.replace('/docs/', '/docs/' + nextVersion);
        }
      } else {
        url = url.replace('/docs/', '/docs/' + nextVersion);
      }
      window.location.href = url;
    }
    // update link url
    const urls = {
      header: document.getElementsByClassName('header__link'),
      footer: document.getElementsByClassName('footer__link'),
      lp: document.getElementsByClassName('mdl-link'),
    };
    for (const links in urls) {
      if (urls[links] != undefined || urls[links] != null) {
        for (var link of urls[links]) {
          if (link.href.includes(beforeVersion)) {
            link.href = link.href.replace('/' + beforeVersion, '');
          }
          if (!elem.className.includes('latest') && link.href.includes('/docs')) {
            link.href = link.href.replace('/docs', '/docs/' + elem.text);
          }
        }
      }
    }
  }
  document.getElementById('version_details').removeAttribute('open');
}

// get github star
const getGitHubStar = async () => {
  const res = await fetch('https://api.github.com/repos/vdaas/vald', {
    method: "GET",
    mode: "cors",
  });
  const json = await res.json()
  return json
}
