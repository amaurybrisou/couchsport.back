import axios from "axios";

const instance = axios.create({
  baseURL: process.env.baseURL,
  withCredentials: true
});

export default instance;
