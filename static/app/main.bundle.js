/******/ (function() { // webpackBootstrap
/******/ 	var __webpack_modules__ = ({

/***/ "./webapp/app/app.tsx":
/*!****************************!*\
  !*** ./webapp/app/app.tsx ***!
  \****************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "App": function() { return /* binding */ App; }
/* harmony export */ });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./node_modules/react/index.js");
/* harmony import */ var _modules_home_home__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./modules/home/home */ "./webapp/app/modules/home/home.tsx");
/* harmony import */ var _elastic_eui_dist_eui_theme_dark_css__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @elastic/eui/dist/eui_theme_dark.css */ "./node_modules/@elastic/eui/dist/eui_theme_dark.css");
/* harmony import */ var react_router_dom__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! react-router-dom */ "./node_modules/react-router-dom/dist/index.js");
/* harmony import */ var react_router_dom__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! react-router-dom */ "./node_modules/react-router/dist/index.js");




const App = () => {
    return (react__WEBPACK_IMPORTED_MODULE_0__.createElement(react_router_dom__WEBPACK_IMPORTED_MODULE_3__.BrowserRouter, null,
        react__WEBPACK_IMPORTED_MODULE_0__.createElement(react_router_dom__WEBPACK_IMPORTED_MODULE_4__.Routes, null,
            react__WEBPACK_IMPORTED_MODULE_0__.createElement(react_router_dom__WEBPACK_IMPORTED_MODULE_4__.Route, { path: "/", element: react__WEBPACK_IMPORTED_MODULE_0__.createElement(_modules_home_home__WEBPACK_IMPORTED_MODULE_1__["default"], null) }))));
};
/* harmony default export */ __webpack_exports__["default"] = (App);


/***/ }),

/***/ "./webapp/app/config/constants.ts":
/*!****************************************!*\
  !*** ./webapp/app/config/constants.ts ***!
  \****************************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "AUTHORITIES": function() { return /* binding */ AUTHORITIES; },
/* harmony export */   "messages": function() { return /* binding */ messages; },
/* harmony export */   "apiUrl": function() { return /* binding */ apiUrl; },
/* harmony export */   "ALBUM_PERMISSIONS": function() { return /* binding */ ALBUM_PERMISSIONS; }
/* harmony export */ });
const AUTHORITIES = {
    ADMIN: 'admin',
    USER: 'user',
    EDITOR: 'editor',
};
const messages = {
    DATA_ERROR_ALERT: 'Internal Error',
};
const apiUrl = 'api/v1';
const ALBUM_PERMISSIONS = {
    READ: 'album.read',
    WRITE: 'album.write',
    EDIT: 'album.edit',
    DELETE: 'album.delete',
};
const isoDateFormat = /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d*)?$/;


/***/ }),

/***/ "./webapp/app/config/error-middleware.ts":
/*!***********************************************!*\
  !*** ./webapp/app/config/error-middleware.ts ***!
  \***********************************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
const getErrorMessage = errorData => {
    let message = errorData.message;
    if (errorData.fieldErrors) {
        errorData.fieldErrors.forEach(fErr => {
            message += `\nfield: ${fErr.field},  Object: ${fErr.objectName}, message: ${fErr.message}\n`;
        });
    }
    return message;
};
/* harmony default export */ __webpack_exports__["default"] = (() => next => action => {
    /**
     *
     * The error middleware serves to log error messages from dispatch
     * It need not run in production
     */
    if (true) {
        const { error } = action;
        if (error) {
            console.error(`${action.type} caught at middleware with reason: ${JSON.stringify(error.message)}.`);
            if (error.response && error.response.data) {
                const message = getErrorMessage(error.response.data);
                console.error(`Actual cause: ${message}`);
            }
        }
    }
    // Dispatch initial action
    return next(action);
});


/***/ }),

/***/ "./webapp/app/config/logger-middleware.ts":
/*!************************************************!*\
  !*** ./webapp/app/config/logger-middleware.ts ***!
  \************************************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* eslint no-console: off */
/* harmony default export */ __webpack_exports__["default"] = (() => next => action => {
    if (true) {
        const { type, payload, meta, error } = action;
        console.groupCollapsed(type);
        console.log('Payload:', payload);
        if (error) {
            console.log('Error:', error);
        }
        console.log('Meta:', meta);
        console.groupEnd();
    }
    return next(action);
});


/***/ }),

/***/ "./webapp/app/config/notification-middleware.ts":
/*!******************************************************!*\
  !*** ./webapp/app/config/notification-middleware.ts ***!
  \******************************************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var react_jhipster__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react-jhipster */ "./node_modules/react-jhipster/lib/index.js");
/* harmony import */ var react_jhipster__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(react_jhipster__WEBPACK_IMPORTED_MODULE_0__);
/* harmony import */ var react_toastify__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! react-toastify */ "./node_modules/react-toastify/dist/react-toastify.esm.js");
/* harmony import */ var app_shared_reducers_reducer_utils__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! app/shared/reducers/reducer.utils */ "./webapp/app/shared/reducers/reducer.utils.ts");



const addErrorAlert = (message, key, data) => {
    key = key ? key : message;
    react_toastify__WEBPACK_IMPORTED_MODULE_1__.toast.error((0,react_jhipster__WEBPACK_IMPORTED_MODULE_0__.translate)(key, data));
};
/* harmony default export */ __webpack_exports__["default"] = (() => next => action => {
    const { error, payload } = action;
    /**
     *
     * The notification middleware serves to add success and error notifications
     */
    if ((0,app_shared_reducers_reducer_utils__WEBPACK_IMPORTED_MODULE_2__.isFulfilledAction)(action) && payload && payload.headers) {
        const headers = payload === null || payload === void 0 ? void 0 : payload.headers;
        let alert = null;
        let alertParams = null;
        headers &&
            Object.entries(headers).forEach(([k, v]) => {
                if (k.toLowerCase().endsWith('app-alert')) {
                    alert = v;
                }
                else if (k.toLowerCase().endsWith('app-params')) {
                    alertParams = decodeURIComponent(v.replace(/\+/g, ' '));
                }
            });
        if (alert) {
            const alertParam = alertParams;
            react_toastify__WEBPACK_IMPORTED_MODULE_1__.toast.success((0,react_jhipster__WEBPACK_IMPORTED_MODULE_0__.translate)(alert, { param: alertParam }));
        }
    }
    if ((0,app_shared_reducers_reducer_utils__WEBPACK_IMPORTED_MODULE_2__.isRejectedAction)(action) && error && error.isAxiosError) {
        if (error.response) {
            const response = error.response;
            const data = response.data;
            switch (response.status) {
                // connection refused, server not reachable
                case 0:
                    addErrorAlert('Server not reachable', 'error.server.not.reachable');
                    break;
                case 400: {
                    let errorHeader = null;
                    let entityKey = null;
                    (response === null || response === void 0 ? void 0 : response.headers) &&
                        Object.entries(response.headers).forEach(([k, v]) => {
                            if (k.toLowerCase().endsWith('app-error')) {
                                errorHeader = v;
                            }
                            else if (k.toLowerCase().endsWith('app-params')) {
                                entityKey = v;
                            }
                        });
                    if (errorHeader) {
                        addErrorAlert(errorHeader, errorHeader, { entityKey });
                    }
                    else if (data === null || data === void 0 ? void 0 : data.fieldErrors) {
                        const fieldErrors = data.fieldErrors;
                        for (const fieldError of fieldErrors) {
                            if (['Min', 'Max', 'DecimalMin', 'DecimalMax'].includes(fieldError.message)) {
                                fieldError.message = 'Size';
                            }
                            // convert 'something[14].other[4].id' to 'something[].other[].id' so translations can be written to it
                            const convertedField = fieldError.field.replace(/\[\d*\]/g, '[]');
                            const fieldName = "${fieldError.objectName}.${convertedField}";
                            addErrorAlert(`Error on field "${convertedField}"`, `error.${fieldError.message}`, { fieldName });
                        }
                    }
                    else if (typeof data === 'string' && data !== '') {
                        addErrorAlert(data);
                    }
                    else {
                        react_toastify__WEBPACK_IMPORTED_MODULE_1__.toast.error((data === null || data === void 0 ? void 0 : data.message) || (data === null || data === void 0 ? void 0 : data.error) || (data === null || data === void 0 ? void 0 : data.title) || 'Unknown error!');
                    }
                    break;
                }
                case 404:
                    addErrorAlert('Not found', 'error.url.not.found');
                    break;
                default:
                    if (typeof data === 'string' && data !== '') {
                        addErrorAlert(data);
                    }
                    else {
                        react_toastify__WEBPACK_IMPORTED_MODULE_1__.toast.error((data === null || data === void 0 ? void 0 : data.message) || (data === null || data === void 0 ? void 0 : data.error) || (data === null || data === void 0 ? void 0 : data.title) || 'Unknown error!');
                    }
            }
        }
    }
    else if (error) {
        react_toastify__WEBPACK_IMPORTED_MODULE_1__.toast.error(error.message || 'Unknown error!');
    }
    return next(action);
});


/***/ }),

/***/ "./webapp/app/config/store.ts":
/*!************************************!*\
  !*** ./webapp/app/config/store.ts ***!
  \************************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "useAppSelector": function() { return /* binding */ useAppSelector; },
/* harmony export */   "useAppDispatch": function() { return /* binding */ useAppDispatch; }
/* harmony export */ });
/* harmony import */ var _reduxjs_toolkit__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! @reduxjs/toolkit */ "./node_modules/@reduxjs/toolkit/dist/redux-toolkit.esm.js");
/* harmony import */ var react_redux__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react-redux */ "./node_modules/react-redux/es/index.js");
/* harmony import */ var app_shared_reducers__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! app/shared/reducers */ "./webapp/app/shared/reducers/index.ts");
/* harmony import */ var _error_middleware__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./error-middleware */ "./webapp/app/config/error-middleware.ts");
/* harmony import */ var _notification_middleware__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./notification-middleware */ "./webapp/app/config/notification-middleware.ts");
/* harmony import */ var _logger_middleware__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ./logger-middleware */ "./webapp/app/config/logger-middleware.ts");
/* harmony import */ var react_redux_loading_bar__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! react-redux-loading-bar */ "./node_modules/react-redux-loading-bar/build/index.js");







const store = (0,_reduxjs_toolkit__WEBPACK_IMPORTED_MODULE_6__.configureStore)({
    reducer: app_shared_reducers__WEBPACK_IMPORTED_MODULE_1__["default"],
    middleware: getDefaultMiddleware => getDefaultMiddleware({
        serializableCheck: {
            // Ignore these field paths in all actions
            ignoredActionPaths: ['payload.config', 'payload.request', 'error', 'meta.arg'],
        },
    }).concat(_error_middleware__WEBPACK_IMPORTED_MODULE_2__["default"], _notification_middleware__WEBPACK_IMPORTED_MODULE_3__["default"], (0,react_redux_loading_bar__WEBPACK_IMPORTED_MODULE_5__.loadingBarMiddleware)(), _logger_middleware__WEBPACK_IMPORTED_MODULE_4__["default"]),
});
const getStore = () => store;
const useAppSelector = react_redux__WEBPACK_IMPORTED_MODULE_0__.useSelector;
const useAppDispatch = () => (0,react_redux__WEBPACK_IMPORTED_MODULE_0__.useDispatch)();
/* harmony default export */ __webpack_exports__["default"] = (getStore);


/***/ }),

/***/ "./webapp/app/index.tsx":
/*!******************************!*\
  !*** ./webapp/app/index.tsx ***!
  \******************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./node_modules/react/index.js");
/* harmony import */ var react_dom__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! react-dom */ "./node_modules/react-dom/index.js");
/* harmony import */ var _app__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./app */ "./webapp/app/app.tsx");
/* harmony import */ var _config_store__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./config/store */ "./webapp/app/config/store.ts");
/* harmony import */ var react_redux__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! react-redux */ "./node_modules/react-redux/es/index.js");
/* harmony import */ var _elastic_eui_dist_eui_theme_dark_css__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! @elastic/eui/dist/eui_theme_dark.css */ "./node_modules/@elastic/eui/dist/eui_theme_dark.css");
/* harmony import */ var _elastic_eui__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! @elastic/eui */ "./node_modules/@elastic/eui/es/components/provider/provider.js");







const store = (0,_config_store__WEBPACK_IMPORTED_MODULE_3__["default"])();
const rootEl = document.getElementById('root');
const render = Component => 
// eslint-disable-next-line react/no-render-return-value
react_dom__WEBPACK_IMPORTED_MODULE_1__.render(react__WEBPACK_IMPORTED_MODULE_0__.createElement(_elastic_eui__WEBPACK_IMPORTED_MODULE_6__.EuiProvider, { colorMode: 'dark' },
    react__WEBPACK_IMPORTED_MODULE_0__.createElement(react_redux__WEBPACK_IMPORTED_MODULE_4__.Provider, { store: store },
        react__WEBPACK_IMPORTED_MODULE_0__.createElement("div", null,
            react__WEBPACK_IMPORTED_MODULE_0__.createElement(Component, null)))), rootEl);
render(_app__WEBPACK_IMPORTED_MODULE_2__["default"]);


/***/ }),

/***/ "./webapp/app/modules/home/home.tsx":
/*!******************************************!*\
  !*** ./webapp/app/modules/home/home.tsx ***!
  \******************************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "Home": function() { return /* binding */ Home; }
/* harmony export */ });
/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react */ "./node_modules/react/index.js");
/* harmony import */ var _elastic_eui__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @elastic/eui */ "./node_modules/@elastic/eui/es/components/page/page.js");
/* harmony import */ var _elastic_eui__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @elastic/eui */ "./node_modules/@elastic/eui/es/components/page/page_sidebar/page_sidebar.js");
/* harmony import */ var _elastic_eui__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! @elastic/eui */ "./node_modules/@elastic/eui/es/components/page/page_body/page_body.js");
/* harmony import */ var _elastic_eui__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! @elastic/eui */ "./node_modules/@elastic/eui/es/components/page/page_section/page_section.js");


const Home = () => {
    return (react__WEBPACK_IMPORTED_MODULE_0__.createElement(_elastic_eui__WEBPACK_IMPORTED_MODULE_1__.EuiPage, null,
        react__WEBPACK_IMPORTED_MODULE_0__.createElement(_elastic_eui__WEBPACK_IMPORTED_MODULE_2__.EuiPageSidebar, { paddingSize: "l" }),
        react__WEBPACK_IMPORTED_MODULE_0__.createElement(_elastic_eui__WEBPACK_IMPORTED_MODULE_3__.EuiPageBody, { paddingSize: "none", panelled: "true" },
            react__WEBPACK_IMPORTED_MODULE_0__.createElement(_elastic_eui__WEBPACK_IMPORTED_MODULE_4__.EuiPageSection, null))));
};
/* harmony default export */ __webpack_exports__["default"] = (Home);


/***/ }),

/***/ "./webapp/app/shared/reducers/album-management.ts":
/*!********************************************************!*\
  !*** ./webapp/app/shared/reducers/album-management.ts ***!
  \********************************************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "getAlbums": function() { return /* binding */ getAlbums; },
/* harmony export */   "AlbumManagementSlice": function() { return /* binding */ AlbumManagementSlice; },
/* harmony export */   "reset": function() { return /* binding */ reset; }
/* harmony export */ });
/* harmony import */ var tslib__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! tslib */ "./node_modules/tslib/tslib.es6.js");
/* harmony import */ var axios__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! axios */ "./node_modules/axios/index.js");
/* harmony import */ var axios__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(axios__WEBPACK_IMPORTED_MODULE_0__);
/* harmony import */ var _reduxjs_toolkit__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @reduxjs/toolkit */ "./node_modules/@reduxjs/toolkit/dist/redux-toolkit.esm.js");
/* harmony import */ var app_config_constants__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! app/config/constants */ "./webapp/app/config/constants.ts");




const DEFAULT_PAGE_SIZE = 20;
const initialState = {
    loading: false,
    errorMessage: null,
    albums: [],
    count: 0,
    offset: 0,
    limit: DEFAULT_PAGE_SIZE,
};
const getAlbums = (0,_reduxjs_toolkit__WEBPACK_IMPORTED_MODULE_2__.createAsyncThunk)('albumManagement/fetch_albums', (offset, limit) => (0,tslib__WEBPACK_IMPORTED_MODULE_3__.__awaiter)(void 0, void 0, void 0, function* () {
    const requestUrl = `${app_config_constants__WEBPACK_IMPORTED_MODULE_1__.apiUrl}/albums?offset=${offset}&limit=${limit}`;
    return axios__WEBPACK_IMPORTED_MODULE_0___default().get(requestUrl);
}));
const AlbumManagementSlice = (0,_reduxjs_toolkit__WEBPACK_IMPORTED_MODULE_2__.createSlice)({
    name: 'albumManagement',
    initialState: initialState,
    reducers: {
        reset() {
            return initialState;
        },
    },
    extraReducers(builder) {
        builder
            .addCase(getAlbums.pending, state => {
            state.loading = true;
        })
            .addCase(getAlbums.rejected, (state, action) => (Object.assign(Object.assign({}, state), { loading: false, errorMessage: action.error.message })))
            .addCase(getAlbums.fulfilled, (state, action) => {
            const d = action.payload.data;
            return Object.assign(Object.assign({}, state), { loading: false, albums: d.albums, count: d.count });
        });
    },
});
const { reset } = AlbumManagementSlice.actions;
/* harmony default export */ __webpack_exports__["default"] = (AlbumManagementSlice.reducer);


/***/ }),

/***/ "./webapp/app/shared/reducers/index.ts":
/*!*********************************************!*\
  !*** ./webapp/app/shared/reducers/index.ts ***!
  \*********************************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony import */ var react_redux_loading_bar__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react-redux-loading-bar */ "./node_modules/react-redux-loading-bar/build/index.js");
/* harmony import */ var _user_management__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./user-management */ "./webapp/app/shared/reducers/user-management.ts");
/* harmony import */ var _album_management__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./album-management */ "./webapp/app/shared/reducers/album-management.ts");



const rootReducer = {
    userManagement: _user_management__WEBPACK_IMPORTED_MODULE_1__["default"],
    albumManagement: _album_management__WEBPACK_IMPORTED_MODULE_2__["default"],
    loadingBar: react_redux_loading_bar__WEBPACK_IMPORTED_MODULE_0__.loadingBarReducer,
};
/* harmony default export */ __webpack_exports__["default"] = (rootReducer);


/***/ }),

/***/ "./webapp/app/shared/reducers/reducer.utils.ts":
/*!*****************************************************!*\
  !*** ./webapp/app/shared/reducers/reducer.utils.ts ***!
  \*****************************************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "isRejectedAction": function() { return /* binding */ isRejectedAction; },
/* harmony export */   "isPendingAction": function() { return /* binding */ isPendingAction; },
/* harmony export */   "isFulfilledAction": function() { return /* binding */ isFulfilledAction; },
/* harmony export */   "serializeAxiosError": function() { return /* binding */ serializeAxiosError; },
/* harmony export */   "createEntitySlice": function() { return /* binding */ createEntitySlice; }
/* harmony export */ });
/* harmony import */ var _reduxjs_toolkit__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @reduxjs/toolkit */ "./node_modules/@reduxjs/toolkit/dist/redux-toolkit.esm.js");

/**
 * Check if the async action type is rejected
 */
function isRejectedAction(action) {
    return action.type.endsWith('/rejected');
}
/**
 * Check if the async action type is pending
 */
function isPendingAction(action) {
    return action.type.endsWith('/pending');
}
/**
 * Check if the async action type is completed
 */
function isFulfilledAction(action) {
    return action.type.endsWith('/fulfilled');
}
const commonErrorProperties = ['name', 'message', 'stack', 'code'];
/**
 * serialize function used for async action errors,
 * since the default function from Redux Toolkit strips useful info from axios errors
 */
const serializeAxiosError = (value) => {
    if (typeof value === 'object' && value !== null) {
        if (value.isAxiosError) {
            return value;
        }
        else {
            const simpleError = {};
            for (const property of commonErrorProperties) {
                if (typeof value[property] === 'string') {
                    simpleError[property] = value[property];
                }
            }
            return simpleError;
        }
    }
    return { message: String(value) };
};
/**
 * A wrapper on top of createSlice from Redux Toolkit to extract
 * common reducers and matchers used by entities
 */
const createEntitySlice = ({ name = '', initialState, reducers, extraReducers, skipRejectionHandling, }) => {
    return (0,_reduxjs_toolkit__WEBPACK_IMPORTED_MODULE_0__.createSlice)({
        name,
        initialState,
        reducers: Object.assign({ 
            /**
             * Reset the entity state to initial state
             */
            reset() {
                return initialState;
            } }, reducers),
        extraReducers(builder) {
            extraReducers(builder);
            /*
             * Common rejection logic is handled here.
             * If you want to add your own rejcetion logic, pass `skipRejectionHandling: true`
             * while calling `createEntitySlice`
             * */
            if (!skipRejectionHandling) {
                builder.addMatcher(isRejectedAction, (state, action) => {
                    state.loading = false;
                    state.updating = false;
                    state.updateSuccess = false;
                    state.errorMessage = action.error.message;
                });
            }
        },
    });
};


/***/ }),

/***/ "./webapp/app/shared/reducers/user-management.ts":
/*!*******************************************************!*\
  !*** ./webapp/app/shared/reducers/user-management.ts ***!
  \*******************************************************/
/***/ (function(__unused_webpack_module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export */ __webpack_require__.d(__webpack_exports__, {
/* harmony export */   "getUsers": function() { return /* binding */ getUsers; },
/* harmony export */   "getAccount": function() { return /* binding */ getAccount; },
/* harmony export */   "UserManagementSlice": function() { return /* binding */ UserManagementSlice; },
/* harmony export */   "reset": function() { return /* binding */ reset; }
/* harmony export */ });
/* harmony import */ var tslib__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! tslib */ "./node_modules/tslib/tslib.es6.js");
/* harmony import */ var axios__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! axios */ "./node_modules/axios/index.js");
/* harmony import */ var axios__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(axios__WEBPACK_IMPORTED_MODULE_0__);
/* harmony import */ var _reduxjs_toolkit__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! @reduxjs/toolkit */ "./node_modules/@reduxjs/toolkit/dist/redux-toolkit.esm.js");
/* harmony import */ var app_config_constants__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! app/config/constants */ "./webapp/app/config/constants.ts");




const initialState = {
    loading: false,
    errorMessage: null,
    users: [],
    account: {},
};
// Async Actions
const getUsers = (0,_reduxjs_toolkit__WEBPACK_IMPORTED_MODULE_2__.createAsyncThunk)('userManagement/fetch_users', () => (0,tslib__WEBPACK_IMPORTED_MODULE_3__.__awaiter)(void 0, void 0, void 0, function* () {
    const requestUrl = `${app_config_constants__WEBPACK_IMPORTED_MODULE_1__.apiUrl}/users`;
    return axios__WEBPACK_IMPORTED_MODULE_0___default().get(requestUrl);
}));
const getAccount = (0,_reduxjs_toolkit__WEBPACK_IMPORTED_MODULE_2__.createAsyncThunk)('authentication/get_account', () => (0,tslib__WEBPACK_IMPORTED_MODULE_3__.__awaiter)(void 0, void 0, void 0, function* () {
    return axios__WEBPACK_IMPORTED_MODULE_0___default().get(`${app_config_constants__WEBPACK_IMPORTED_MODULE_1__.apiUrl}/account`);
}));
const UserManagementSlice = (0,_reduxjs_toolkit__WEBPACK_IMPORTED_MODULE_2__.createSlice)({
    name: 'userManagement',
    initialState: initialState,
    reducers: {
        reset() {
            return initialState;
        },
    },
    extraReducers(builder) {
        builder
            .addCase(getUsers.pending, state => {
            state.loading = true;
        })
            .addCase(getUsers.rejected, (state, action) => (Object.assign(Object.assign({}, state), { loading: false, errorMessage: action.error.message })))
            .addCase(getUsers.fulfilled, (state, action) => (Object.assign(Object.assign({}, state), { loading: false, users: action.payload.data })))
            .addCase(getAccount.pending, state => {
            state.loading = true;
        })
            .addCase(getAccount.rejected, (state, action) => (Object.assign(Object.assign({}, state), { loading: false, errorMessage: action.error.message })))
            .addCase(getAccount.fulfilled, (state, action) => (Object.assign(Object.assign({}, state), { loading: false, account: action.payload.data })));
    },
});
const { reset } = UserManagementSlice.actions;
// Reducer
/* harmony default export */ __webpack_exports__["default"] = (UserManagementSlice.reducer);


/***/ }),

/***/ "data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%2710%27 height=%278%27 viewBox=%270 0 10 8%27%3E%3Cpath d=%27M.375 2.625L3.375 5.625M3.375 5.625L8.625.375%27 fill=%27none%27 fill-rule=%27evenodd%27 stroke=%27rgb%28166, 170, 176%29%27 stroke-linecap=%27round%27 stroke-width=%271.5%27 transform=%27translate%28.5 1%29%27/%3E%3C/svg%3E":
/*!*****************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%2710%27 height=%278%27 viewBox=%270 0 10 8%27%3E%3Cpath d=%27M.375 2.625L3.375 5.625M3.375 5.625L8.625.375%27 fill=%27none%27 fill-rule=%27evenodd%27 stroke=%27rgb%28166, 170, 176%29%27 stroke-linecap=%27round%27 stroke-width=%271.5%27 transform=%27translate%28.5 1%29%27/%3E%3C/svg%3E ***!
  \*****************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************/
/***/ (function(module) {

"use strict";
module.exports = "data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%2710%27 height=%278%27 viewBox=%270 0 10 8%27%3E%3Cpath d=%27M.375 2.625L3.375 5.625M3.375 5.625L8.625.375%27 fill=%27none%27 fill-rule=%27evenodd%27 stroke=%27rgb%28166, 170, 176%29%27 stroke-linecap=%27round%27 stroke-width=%271.5%27 transform=%27translate%28.5 1%29%27/%3E%3C/svg%3E";

/***/ }),

/***/ "data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%2710%27 height=%278%27 viewBox=%270 0 10 8%27%3E%3Cpath d=%27M.375 2.625L3.375 5.625M3.375 5.625L8.625.375%27 fill=%27none%27 fill-rule=%27evenodd%27 stroke=%27rgb%2829, 30, 36%29%27 stroke-linecap=%27round%27 stroke-width=%271.5%27 transform=%27translate%28.5 1%29%27/%3E%3C/svg%3E":
/*!**************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%2710%27 height=%278%27 viewBox=%270 0 10 8%27%3E%3Cpath d=%27M.375 2.625L3.375 5.625M3.375 5.625L8.625.375%27 fill=%27none%27 fill-rule=%27evenodd%27 stroke=%27rgb%2829, 30, 36%29%27 stroke-linecap=%27round%27 stroke-width=%271.5%27 transform=%27translate%28.5 1%29%27/%3E%3C/svg%3E ***!
  \**************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************/
/***/ (function(module) {

"use strict";
module.exports = "data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%2710%27 height=%278%27 viewBox=%270 0 10 8%27%3E%3Cpath d=%27M.375 2.625L3.375 5.625M3.375 5.625L8.625.375%27 fill=%27none%27 fill-rule=%27evenodd%27 stroke=%27rgb%2829, 30, 36%29%27 stroke-linecap=%27round%27 stroke-width=%271.5%27 transform=%27translate%28.5 1%29%27/%3E%3C/svg%3E";

/***/ }),

/***/ "data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%276%27 height=%276%27 viewBox=%270 0 6 6%27%3E%3Ccircle cx=%278%27 cy=%2711%27 r=%273%27 fill=%27rgb%28166, 170, 176%29%27 fill-rule=%27evenodd%27 transform=%27translate%28-5 -8%29%27/%3E%3C/svg%3E":
/*!*************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%276%27 height=%276%27 viewBox=%270 0 6 6%27%3E%3Ccircle cx=%278%27 cy=%2711%27 r=%273%27 fill=%27rgb%28166, 170, 176%29%27 fill-rule=%27evenodd%27 transform=%27translate%28-5 -8%29%27/%3E%3C/svg%3E ***!
  \*************************************************************************************************************************************************************************************************************************************************************************************/
/***/ (function(module) {

"use strict";
module.exports = "data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%276%27 height=%276%27 viewBox=%270 0 6 6%27%3E%3Ccircle cx=%278%27 cy=%2711%27 r=%273%27 fill=%27rgb%28166, 170, 176%29%27 fill-rule=%27evenodd%27 transform=%27translate%28-5 -8%29%27/%3E%3C/svg%3E";

/***/ }),

/***/ "data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%276%27 height=%276%27 viewBox=%270 0 6 6%27%3E%3Ccircle cx=%278%27 cy=%2711%27 r=%273%27 fill=%27rgb%2829, 30, 36%29%27 fill-rule=%27evenodd%27 transform=%27translate%28-5 -8%29%27/%3E%3C/svg%3E":
/*!**********************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%276%27 height=%276%27 viewBox=%270 0 6 6%27%3E%3Ccircle cx=%278%27 cy=%2711%27 r=%273%27 fill=%27rgb%2829, 30, 36%29%27 fill-rule=%27evenodd%27 transform=%27translate%28-5 -8%29%27/%3E%3C/svg%3E ***!
  \**********************************************************************************************************************************************************************************************************************************************************************************/
/***/ (function(module) {

"use strict";
module.exports = "data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%276%27 height=%276%27 viewBox=%270 0 6 6%27%3E%3Ccircle cx=%278%27 cy=%2711%27 r=%273%27 fill=%27rgb%2829, 30, 36%29%27 fill-rule=%27evenodd%27 transform=%27translate%28-5 -8%29%27/%3E%3C/svg%3E";

/***/ }),

/***/ "data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%276%27 height=%276%27 viewBox=%270 0 6 6%27%3E%3Crect width=%276%27 height=%276%27 fill=%27rgb%28166, 170, 176%29%27 fill-rule=%27evenodd%27/%3E%3C/svg%3E":
/*!******************************************************************************************************************************************************************************************************************************************!*\
  !*** data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%276%27 height=%276%27 viewBox=%270 0 6 6%27%3E%3Crect width=%276%27 height=%276%27 fill=%27rgb%28166, 170, 176%29%27 fill-rule=%27evenodd%27/%3E%3C/svg%3E ***!
  \******************************************************************************************************************************************************************************************************************************************/
/***/ (function(module) {

"use strict";
module.exports = "data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%276%27 height=%276%27 viewBox=%270 0 6 6%27%3E%3Crect width=%276%27 height=%276%27 fill=%27rgb%28166, 170, 176%29%27 fill-rule=%27evenodd%27/%3E%3C/svg%3E";

/***/ }),

/***/ "data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%276%27 height=%276%27 viewBox=%270 0 6 6%27%3E%3Crect width=%276%27 height=%276%27 fill=%27rgb%2829, 30, 36%29%27 fill-rule=%27evenodd%27/%3E%3C/svg%3E":
/*!***************************************************************************************************************************************************************************************************************************************!*\
  !*** data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%276%27 height=%276%27 viewBox=%270 0 6 6%27%3E%3Crect width=%276%27 height=%276%27 fill=%27rgb%2829, 30, 36%29%27 fill-rule=%27evenodd%27/%3E%3C/svg%3E ***!
  \***************************************************************************************************************************************************************************************************************************************/
/***/ (function(module) {

"use strict";
module.exports = "data:image/svg+xml,%3Csvg xmlns=%27http://www.w3.org/2000/svg%27 width=%276%27 height=%276%27 viewBox=%270 0 6 6%27%3E%3Crect width=%276%27 height=%276%27 fill=%27rgb%2829, 30, 36%29%27 fill-rule=%27evenodd%27/%3E%3C/svg%3E";

/***/ }),

/***/ "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxNiIgaGVpZ2h0PSIxNiI+PHBhdGggZmlsbD0iI0ZGRiIgZD0iTTExLjY5MiA3SDMuNTU2QzMuMjQ5IDcgMyA3LjIyNCAzIDcuNXMuMjQ5LjUuNTU2LjVoOC4xMzZsLTQuMDk2IDQuMDk2YS41LjUgMCAwMC43MDcuNzA3bDQuMjQzLTQuMjQyYy4yNTgtLjI1OS40MDMtLjU4Ny40MzMtLjkyNWEuNDU0LjQ1NCAwIDAwMC0uMjcyIDEuNDk0IDEuNDk0IDAgMDAtLjQzMy0uOTI1TDguMzAzIDIuMTk3YS41LjUgMCAxMC0uNzA3LjcwN0wxMS42OTIgN3oiLz48L3N2Zz4=":
/*!**************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxNiIgaGVpZ2h0PSIxNiI+PHBhdGggZmlsbD0iI0ZGRiIgZD0iTTExLjY5MiA3SDMuNTU2QzMuMjQ5IDcgMyA3LjIyNCAzIDcuNXMuMjQ5LjUuNTU2LjVoOC4xMzZsLTQuMDk2IDQuMDk2YS41LjUgMCAwMC43MDcuNzA3bDQuMjQzLTQuMjQyYy4yNTgtLjI1OS40MDMtLjU4Ny40MzMtLjkyNWEuNDU0LjQ1NCAwIDAwMC0uMjcyIDEuNDk0IDEuNDk0IDAgMDAtLjQzMy0uOTI1TDguMzAzIDIuMTk3YS41LjUgMCAxMC0uNzA3LjcwN0wxMS42OTIgN3oiLz48L3N2Zz4= ***!
  \**************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************/
/***/ (function(module) {

"use strict";
module.exports = "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSIxNiIgaGVpZ2h0PSIxNiI+PHBhdGggZmlsbD0iI0ZGRiIgZD0iTTExLjY5MiA3SDMuNTU2QzMuMjQ5IDcgMyA3LjIyNCAzIDcuNXMuMjQ5LjUuNTU2LjVoOC4xMzZsLTQuMDk2IDQuMDk2YS41LjUgMCAwMC43MDcuNzA3bDQuMjQzLTQuMjQyYy4yNTgtLjI1OS40MDMtLjU4Ny40MzMtLjkyNWEuNDU0LjQ1NCAwIDAwMC0uMjcyIDEuNDk0IDEuNDk0IDAgMDAtLjQzMy0uOTI1TDguMzAzIDIuMTk3YS41LjUgMCAxMC0uNzA3LjcwN0wxMS42OTIgN3oiLz48L3N2Zz4=";

/***/ }),

/***/ "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiB3aWR0aD0iMTYiIGhlaWdodD0iMTYiIHZpZXdCb3g9IjAgMCAxNiAxNiI+ICA8ZGVmcz4gICAgPHBhdGggaWQ9ImRvd25hcnJvd3doaXRlLWEiIGQ9Ik0xMy4wNjg4NTA4LDUuMTU3MjUwMzggTDguMzg0MjM5NzUsOS43NjgyNzQyOCBDOC4xNzA1NDQxNSw5Ljk3ODYxMzA4IDcuODI5OTkyMTQsOS45NzkxNDA5NSA3LjYxNTc2MDI1LDkuNzY4Mjc0MjggTDIuOTMxMTQ5MTUsNS4xNTcyNTAzOCBDMi43MTgxMzU5LDQuOTQ3NTgzMjEgMi4zNzI3NzMxOSw0Ljk0NzU4MzIxIDIuMTU5NzU5OTQsNS4xNTcyNTAzOCBDMS45NDY3NDY2OSw1LjM2NjkxNzU2IDEuOTQ2NzQ2NjksNS43MDY4NTUyMiAyLjE1OTc1OTk0LDUuOTE2NTIyNCBMNi44NDQzNzEwNCwxMC41Mjc1NDYzIEM3LjQ4NTE3NDI0LDExLjE1ODI4MzYgOC41MTY0NDk3OSwxMS4xNTY2ODUxIDkuMTU1NjI4OTYsMTAuNTI3NTQ2MyBMMTMuODQwMjQwMSw1LjkxNjUyMjQgQzE0LjA1MzI1MzMsNS43MDY4NTUyMiAxNC4wNTMyNTMzLDUuMzY2OTE3NTYgMTMuODQwMjQwMSw1LjE1NzI1MDM4IEMxMy42MjcyMjY4LDQuOTQ3NTgzMjEgMTMuMjgxODY0MSw0Ljk0NzU4MzIxIDEzLjA2ODg1MDgsNS4xNTcyNTAzOCBaIi8+ICA8L2RlZnM+ICA8ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPiAgICA8dXNlIGZpbGw9IiNGRkYiIGZpbGwtcnVsZT0ibm9uemVybyIgeGxpbms6aHJlZj0iI2Rvd25hcnJvd3doaXRlLWEiLz4gIDwvZz48L3N2Zz4=":
/*!**********************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************!*\
  !*** data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiB3aWR0aD0iMTYiIGhlaWdodD0iMTYiIHZpZXdCb3g9IjAgMCAxNiAxNiI+ICA8ZGVmcz4gICAgPHBhdGggaWQ9ImRvd25hcnJvd3doaXRlLWEiIGQ9Ik0xMy4wNjg4NTA4LDUuMTU3MjUwMzggTDguMzg0MjM5NzUsOS43NjgyNzQyOCBDOC4xNzA1NDQxNSw5Ljk3ODYxMzA4IDcuODI5OTkyMTQsOS45NzkxNDA5NSA3LjYxNTc2MDI1LDkuNzY4Mjc0MjggTDIuOTMxMTQ5MTUsNS4xNTcyNTAzOCBDMi43MTgxMzU5LDQuOTQ3NTgzMjEgMi4zNzI3NzMxOSw0Ljk0NzU4MzIxIDIuMTU5NzU5OTQsNS4xNTcyNTAzOCBDMS45NDY3NDY2OSw1LjM2NjkxNzU2IDEuOTQ2NzQ2NjksNS43MDY4NTUyMiAyLjE1OTc1OTk0LDUuOTE2NTIyNCBMNi44NDQzNzEwNCwxMC41Mjc1NDYzIEM3LjQ4NTE3NDI0LDExLjE1ODI4MzYgOC41MTY0NDk3OSwxMS4xNTY2ODUxIDkuMTU1NjI4OTYsMTAuNTI3NTQ2MyBMMTMuODQwMjQwMSw1LjkxNjUyMjQgQzE0LjA1MzI1MzMsNS43MDY4NTUyMiAxNC4wNTMyNTMzLDUuMzY2OTE3NTYgMTMuODQwMjQwMSw1LjE1NzI1MDM4IEMxMy42MjcyMjY4LDQuOTQ3NTgzMjEgMTMuMjgxODY0MSw0Ljk0NzU4MzIxIDEzLjA2ODg1MDgsNS4xNTcyNTAzOCBaIi8+ICA8L2RlZnM+ICA8ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPiAgICA8dXNlIGZpbGw9IiNGRkYiIGZpbGwtcnVsZT0ibm9uemVybyIgeGxpbms6aHJlZj0iI2Rvd25hcnJvd3doaXRlLWEiLz4gIDwvZz48L3N2Zz4= ***!
  \**********************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************************/
/***/ (function(module) {

"use strict";
module.exports = "data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiB3aWR0aD0iMTYiIGhlaWdodD0iMTYiIHZpZXdCb3g9IjAgMCAxNiAxNiI+ICA8ZGVmcz4gICAgPHBhdGggaWQ9ImRvd25hcnJvd3doaXRlLWEiIGQ9Ik0xMy4wNjg4NTA4LDUuMTU3MjUwMzggTDguMzg0MjM5NzUsOS43NjgyNzQyOCBDOC4xNzA1NDQxNSw5Ljk3ODYxMzA4IDcuODI5OTkyMTQsOS45NzkxNDA5NSA3LjYxNTc2MDI1LDkuNzY4Mjc0MjggTDIuOTMxMTQ5MTUsNS4xNTcyNTAzOCBDMi43MTgxMzU5LDQuOTQ3NTgzMjEgMi4zNzI3NzMxOSw0Ljk0NzU4MzIxIDIuMTU5NzU5OTQsNS4xNTcyNTAzOCBDMS45NDY3NDY2OSw1LjM2NjkxNzU2IDEuOTQ2NzQ2NjksNS43MDY4NTUyMiAyLjE1OTc1OTk0LDUuOTE2NTIyNCBMNi44NDQzNzEwNCwxMC41Mjc1NDYzIEM3LjQ4NTE3NDI0LDExLjE1ODI4MzYgOC41MTY0NDk3OSwxMS4xNTY2ODUxIDkuMTU1NjI4OTYsMTAuNTI3NTQ2MyBMMTMuODQwMjQwMSw1LjkxNjUyMjQgQzE0LjA1MzI1MzMsNS43MDY4NTUyMiAxNC4wNTMyNTMzLDUuMzY2OTE3NTYgMTMuODQwMjQwMSw1LjE1NzI1MDM4IEMxMy42MjcyMjY4LDQuOTQ3NTgzMjEgMTMuMjgxODY0MSw0Ljk0NzU4MzIxIDEzLjA2ODg1MDgsNS4xNTcyNTAzOCBaIi8+ICA8L2RlZnM+ICA8ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPiAgICA8dXNlIGZpbGw9IiNGRkYiIGZpbGwtcnVsZT0ibm9uemVybyIgeGxpbms6aHJlZj0iI2Rvd25hcnJvd3doaXRlLWEiLz4gIDwvZz48L3N2Zz4=";

/***/ }),

/***/ "?5580":
/*!**************************************!*\
  !*** ./terminal-highlight (ignored) ***!
  \**************************************/
/***/ (function() {

/* (ignored) */

/***/ }),

/***/ "?03fb":
/*!********************!*\
  !*** fs (ignored) ***!
  \********************/
/***/ (function() {

/* (ignored) */

/***/ }),

/***/ "?6197":
/*!**********************!*\
  !*** path (ignored) ***!
  \**********************/
/***/ (function() {

/* (ignored) */

/***/ }),

/***/ "?b8cb":
/*!*******************************!*\
  !*** source-map-js (ignored) ***!
  \*******************************/
/***/ (function() {

/* (ignored) */

/***/ }),

/***/ "?c717":
/*!*********************!*\
  !*** url (ignored) ***!
  \*********************/
/***/ (function() {

/* (ignored) */

/***/ })

/******/ 	});
/************************************************************************/
/******/ 	// The module cache
/******/ 	var __webpack_module_cache__ = {};
/******/ 	
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/ 		// Check if module is in cache
/******/ 		var cachedModule = __webpack_module_cache__[moduleId];
/******/ 		if (cachedModule !== undefined) {
/******/ 			return cachedModule.exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = __webpack_module_cache__[moduleId] = {
/******/ 			id: moduleId,
/******/ 			loaded: false,
/******/ 			exports: {}
/******/ 		};
/******/ 	
/******/ 		// Execute the module function
/******/ 		__webpack_modules__[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/ 	
/******/ 		// Flag the module as loaded
/******/ 		module.loaded = true;
/******/ 	
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/ 	
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = __webpack_modules__;
/******/ 	
/************************************************************************/
/******/ 	/* webpack/runtime/chunk loaded */
/******/ 	!function() {
/******/ 		var deferred = [];
/******/ 		__webpack_require__.O = function(result, chunkIds, fn, priority) {
/******/ 			if(chunkIds) {
/******/ 				priority = priority || 0;
/******/ 				for(var i = deferred.length; i > 0 && deferred[i - 1][2] > priority; i--) deferred[i] = deferred[i - 1];
/******/ 				deferred[i] = [chunkIds, fn, priority];
/******/ 				return;
/******/ 			}
/******/ 			var notFulfilled = Infinity;
/******/ 			for (var i = 0; i < deferred.length; i++) {
/******/ 				var chunkIds = deferred[i][0];
/******/ 				var fn = deferred[i][1];
/******/ 				var priority = deferred[i][2];
/******/ 				var fulfilled = true;
/******/ 				for (var j = 0; j < chunkIds.length; j++) {
/******/ 					if ((priority & 1 === 0 || notFulfilled >= priority) && Object.keys(__webpack_require__.O).every(function(key) { return __webpack_require__.O[key](chunkIds[j]); })) {
/******/ 						chunkIds.splice(j--, 1);
/******/ 					} else {
/******/ 						fulfilled = false;
/******/ 						if(priority < notFulfilled) notFulfilled = priority;
/******/ 					}
/******/ 				}
/******/ 				if(fulfilled) {
/******/ 					deferred.splice(i--, 1)
/******/ 					var r = fn();
/******/ 					if (r !== undefined) result = r;
/******/ 				}
/******/ 			}
/******/ 			return result;
/******/ 		};
/******/ 	}();
/******/ 	
/******/ 	/* webpack/runtime/compat get default export */
/******/ 	!function() {
/******/ 		// getDefaultExport function for compatibility with non-harmony modules
/******/ 		__webpack_require__.n = function(module) {
/******/ 			var getter = module && module.__esModule ?
/******/ 				function() { return module['default']; } :
/******/ 				function() { return module; };
/******/ 			__webpack_require__.d(getter, { a: getter });
/******/ 			return getter;
/******/ 		};
/******/ 	}();
/******/ 	
/******/ 	/* webpack/runtime/create fake namespace object */
/******/ 	!function() {
/******/ 		var getProto = Object.getPrototypeOf ? function(obj) { return Object.getPrototypeOf(obj); } : function(obj) { return obj.__proto__; };
/******/ 		var leafPrototypes;
/******/ 		// create a fake namespace object
/******/ 		// mode & 1: value is a module id, require it
/******/ 		// mode & 2: merge all properties of value into the ns
/******/ 		// mode & 4: return value when already ns object
/******/ 		// mode & 16: return value when it's Promise-like
/******/ 		// mode & 8|1: behave like require
/******/ 		__webpack_require__.t = function(value, mode) {
/******/ 			if(mode & 1) value = this(value);
/******/ 			if(mode & 8) return value;
/******/ 			if(typeof value === 'object' && value) {
/******/ 				if((mode & 4) && value.__esModule) return value;
/******/ 				if((mode & 16) && typeof value.then === 'function') return value;
/******/ 			}
/******/ 			var ns = Object.create(null);
/******/ 			__webpack_require__.r(ns);
/******/ 			var def = {};
/******/ 			leafPrototypes = leafPrototypes || [null, getProto({}), getProto([]), getProto(getProto)];
/******/ 			for(var current = mode & 2 && value; typeof current == 'object' && !~leafPrototypes.indexOf(current); current = getProto(current)) {
/******/ 				Object.getOwnPropertyNames(current).forEach(function(key) { def[key] = function() { return value[key]; }; });
/******/ 			}
/******/ 			def['default'] = function() { return value; };
/******/ 			__webpack_require__.d(ns, def);
/******/ 			return ns;
/******/ 		};
/******/ 	}();
/******/ 	
/******/ 	/* webpack/runtime/define property getters */
/******/ 	!function() {
/******/ 		// define getter functions for harmony exports
/******/ 		__webpack_require__.d = function(exports, definition) {
/******/ 			for(var key in definition) {
/******/ 				if(__webpack_require__.o(definition, key) && !__webpack_require__.o(exports, key)) {
/******/ 					Object.defineProperty(exports, key, { enumerable: true, get: definition[key] });
/******/ 				}
/******/ 			}
/******/ 		};
/******/ 	}();
/******/ 	
/******/ 	/* webpack/runtime/global */
/******/ 	!function() {
/******/ 		__webpack_require__.g = (function() {
/******/ 			if (typeof globalThis === 'object') return globalThis;
/******/ 			try {
/******/ 				return this || new Function('return this')();
/******/ 			} catch (e) {
/******/ 				if (typeof window === 'object') return window;
/******/ 			}
/******/ 		})();
/******/ 	}();
/******/ 	
/******/ 	/* webpack/runtime/hasOwnProperty shorthand */
/******/ 	!function() {
/******/ 		__webpack_require__.o = function(obj, prop) { return Object.prototype.hasOwnProperty.call(obj, prop); }
/******/ 	}();
/******/ 	
/******/ 	/* webpack/runtime/make namespace object */
/******/ 	!function() {
/******/ 		// define __esModule on exports
/******/ 		__webpack_require__.r = function(exports) {
/******/ 			if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 				Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 			}
/******/ 			Object.defineProperty(exports, '__esModule', { value: true });
/******/ 		};
/******/ 	}();
/******/ 	
/******/ 	/* webpack/runtime/node module decorator */
/******/ 	!function() {
/******/ 		__webpack_require__.nmd = function(module) {
/******/ 			module.paths = [];
/******/ 			if (!module.children) module.children = [];
/******/ 			return module;
/******/ 		};
/******/ 	}();
/******/ 	
/******/ 	/* webpack/runtime/jsonp chunk loading */
/******/ 	!function() {
/******/ 		__webpack_require__.b = document.baseURI || self.location.href;
/******/ 		
/******/ 		// object to store loaded and loading chunks
/******/ 		// undefined = chunk not loaded, null = chunk preloaded/prefetched
/******/ 		// [resolve, reject, Promise] = chunk loading, 0 = chunk loaded
/******/ 		var installedChunks = {
/******/ 			"main": 0
/******/ 		};
/******/ 		
/******/ 		// no chunk on demand loading
/******/ 		
/******/ 		// no prefetching
/******/ 		
/******/ 		// no preloaded
/******/ 		
/******/ 		// no HMR
/******/ 		
/******/ 		// no HMR manifest
/******/ 		
/******/ 		__webpack_require__.O.j = function(chunkId) { return installedChunks[chunkId] === 0; };
/******/ 		
/******/ 		// install a JSONP callback for chunk loading
/******/ 		var webpackJsonpCallback = function(parentChunkLoadingFunction, data) {
/******/ 			var chunkIds = data[0];
/******/ 			var moreModules = data[1];
/******/ 			var runtime = data[2];
/******/ 			// add "moreModules" to the modules object,
/******/ 			// then flag all "chunkIds" as loaded and fire callback
/******/ 			var moduleId, chunkId, i = 0;
/******/ 			if(chunkIds.some(function(id) { return installedChunks[id] !== 0; })) {
/******/ 				for(moduleId in moreModules) {
/******/ 					if(__webpack_require__.o(moreModules, moduleId)) {
/******/ 						__webpack_require__.m[moduleId] = moreModules[moduleId];
/******/ 					}
/******/ 				}
/******/ 				if(runtime) var result = runtime(__webpack_require__);
/******/ 			}
/******/ 			if(parentChunkLoadingFunction) parentChunkLoadingFunction(data);
/******/ 			for(;i < chunkIds.length; i++) {
/******/ 				chunkId = chunkIds[i];
/******/ 				if(__webpack_require__.o(installedChunks, chunkId) && installedChunks[chunkId]) {
/******/ 					installedChunks[chunkId][0]();
/******/ 				}
/******/ 				installedChunks[chunkId] = 0;
/******/ 			}
/******/ 			return __webpack_require__.O(result);
/******/ 		}
/******/ 		
/******/ 		var chunkLoadingGlobal = self["webpackChunkmy_app"] = self["webpackChunkmy_app"] || [];
/******/ 		chunkLoadingGlobal.forEach(webpackJsonpCallback.bind(null, 0));
/******/ 		chunkLoadingGlobal.push = webpackJsonpCallback.bind(null, chunkLoadingGlobal.push.bind(chunkLoadingGlobal));
/******/ 	}();
/******/ 	
/************************************************************************/
/******/ 	
/******/ 	// startup
/******/ 	// Load entry module and return exports
/******/ 	// This entry module depends on other loaded chunks and execution need to be delayed
/******/ 	var __webpack_exports__ = __webpack_require__.O(undefined, ["vendors"], function() { return __webpack_require__("./webapp/app/index.tsx"); })
/******/ 	__webpack_exports__ = __webpack_require__.O(__webpack_exports__);
/******/ 	
/******/ })()
;
//# sourceMappingURL=main.bundle.js.map