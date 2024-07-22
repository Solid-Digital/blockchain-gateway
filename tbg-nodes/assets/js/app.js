// We need to import the CSS so that webpack will load it.
// The MiniCssExtractPlugin is used to separate it out into
// its own CSS file.
import css from "../css/app.scss"

// webpack automatically bundles all modules in your
// entry points. Those entry points can be configured
// in "webpack.config.js".
//
// Import dependencies
//
import "phoenix_html"
import "../node_modules/@webcomponents/custom-elements/custom-elements.min.js"
import 'alpinejs';

// Import local files
//
// Local files can be imported directly using relative paths, for example:
// import socket from "./socket"
import "./password"
import "./network-detail"

import {Socket} from "phoenix"
import LiveSocket from "phoenix_live_view"


// js for app.html.eex
let csrfToken = document.querySelector("meta[name='csrf-token']").getAttribute("content");
let liveSocket = new LiveSocket("/live", Socket, {params: {_csrf_token: csrfToken}});
liveSocket.connect();
