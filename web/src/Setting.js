// Copyright 2021 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import {message} from "antd";
import React from "react";
import {isMobile as isMobileDevice} from "react-device-detect";
import "./i18n";
import i18next from "i18next";
import copy from "copy-to-clipboard";
import {authConfig} from "./auth/Auth";
import {Helmet} from "react-helmet";

export let ServerUrl = "";

export function initServerUrl() {
  const hostname = window.location.hostname;
  if (hostname === "localhost") {
    ServerUrl = `http://${hostname}:8000`;
  }
}

function isLocalhost() {
  const hostname = window.location.hostname;
  return hostname === "localhost";
}

export function isProviderVisible(provider) {
  if (provider.type !== "GitHub") {
    return true;
  }

  if (isLocalhost()) {
    return provider.name.includes("localhost");
  } else {
    return !provider.name.includes("localhost");
  }
}

export function parseJson(s) {
  if (s === "") {
    return null;
  } else {
    return JSON.parse(s);
  }
}

export function myParseInt(i) {
  const res = parseInt(i);
  return isNaN(res) ? 0 : res;
}

export function openLink(link) {
  // this.props.history.push(link);
  const w = window.open('about:blank');
  w.location.href = link;
}

export function goToLink(link) {
  window.location.href = link;
}

export function goToLinkSoft(ths, link) {
  ths.props.history.push(link);
}

export function showMessage(type, text) {
  if (type === "") {
    return;
  } else if (type === "success") {
    message.success(text);
  } else if (type === "error") {
    message.error(text);
  }
}

export function isAdminUser(account) {
  if (account === undefined || account === null) {
    return false;
  }
  return account.owner === "built-in" || account.isGlobalAdmin === true;
}

export function deepCopy(obj) {
  return Object.assign({}, obj);
}

export function addRow(array, row) {
  return [...array, row];
}

export function prependRow(array, row) {
  return [row, ...array];
}

export function deleteRow(array, i) {
  // return array = array.slice(0, i).concat(array.slice(i + 1));
  return [...array.slice(0, i), ...array.slice(i + 1)];
}

export function swapRow(array, i, j) {
  return [...array.slice(0, i), array[j], ...array.slice(i + 1, j), array[i], ...array.slice(j + 1)];
}

export function isMobile() {
  // return getIsMobileView();
  return isMobileDevice;
}

export function getFormattedDate(date) {
  if (date === undefined) {
    return null;
  }

  date = date.replace('T', ' ');
  date = date.replace('+08:00', ' ');
  return date;
}

export function getFormattedDateShort(date) {
  return date.slice(0, 10);
}

export function getShortName(s) {
  return s.split('/').slice(-1)[0];
}

export function getShortText(s, maxLength=35) {
  if (s.length > maxLength) {
    return `${s.slice(0, maxLength)}...`;
  } else {
    return s;
  }
}

function getRandomInt(s) {
  let hash = 0;
  if (s.length !== 0) {
    for (let i = 0; i < s.length; i ++) {
      let char = s.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash;
    }
  }

  return hash;
}

export function getAvatarColor(s) {
  const colorList = ['#f56a00', '#7265e6', '#ffbf00', '#00a2ae'];
  let random = getRandomInt(s);
  if (random < 0) {
    random = -random;
  }
  return colorList[random % 4];
}

export function setLanguage() {
  let language = localStorage.getItem('language');
  if (language === undefined) {
    language = "en"
  }
  i18next.changeLanguage(language)
}

export function changeLanguage(language) {
  localStorage.setItem("language", language)
  i18next.changeLanguage(language)
  window.location.reload(true);
}

export function getClickable(text) {
  return (
    // eslint-disable-next-line jsx-a11y/anchor-is-valid
    <a onClick={() => {
      copy(text);
      showMessage("success", `Copied to clipboard`);
    }}>
      {text}
    </a>
  )
}

export function getProviderLogo(provider) {
  const idp = provider.type.toLowerCase();
  const url = `https://cdn.jsdelivr.net/gh/casbin/static/img/social_${idp}.png`;
  return (
    <img width={30} height={30} src={url} alt={idp} />
  )
}

export function renderLogo(application) {
  if (application === null) {
    return null;
  }

  if (application.homepageUrl !== "") {
    return (
      <a target="_blank" rel="noreferrer" href={application.homepageUrl}>
        <img width={250} src={application.logo} alt={application.displayName} style={{marginBottom: '30px'}}/>
      </a>
    )
  } else {
    return (
      <img width={250} src={application.logo} alt={application.displayName} style={{marginBottom: '30px'}}/>
    );
  }
}

export function goToLogin(ths, application) {
  if (application === null) {
    return;
  }

  if (authConfig.appName === application.name) {
    goToLinkSoft(ths, "/login");
  } else {
    goToLink(`${application.homepageUrl}/login`);
  }
}

export function goToSignup(ths, application) {
  if (application === null) {
    return;
  }

  if (authConfig.appName === application.name) {
    goToLinkSoft(ths, "/signup");
  } else {
    if (application.signupUrl === "") {
      goToLinkSoft(ths, `/signup/${application.name}`);
    } else {
      goToLink(application.signupUrl);
    }
  }
}

export function goToForget(ths, application) {
  if (application === null) {
    return;
  }

  if (authConfig.appName === application.name) {
    goToLinkSoft(ths, "/forget");
  } else {
    if (application.signupUrl === "") {
      goToLinkSoft(ths, `/forget/${application.name}`);
    } else {
      goToLink(application.forgetUrl);
    }
  }
}

export function renderHelmet(application) {
  if (application === undefined || application === null || application.organizationObj === undefined || application.organizationObj === null ||application.organizationObj === "") {
    return null;
  }

  return (
    <Helmet>
      <title>{application.organizationObj.displayName}</title>
      <link rel="icon" href={application.organizationObj.favicon} />
    </Helmet>
  )
}
