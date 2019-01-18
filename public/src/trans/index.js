import Vue from "vue";
import VueI18n from "vue-i18n";

Vue.use(VueI18n);

import fr from "./fr";
import en from "./en";

const messages = {
  en: en,
  fr: fr
};

export const i18n = new VueI18n({
  locale: "fr", // set locale
  fallbackLocale: "en", // set fallback locale
  messages // set locale messages
});

export const defaultLocale = "fr";
