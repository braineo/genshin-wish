import axios from 'axios';

const client = axios.create({
  baseURL: '/api/v1',
});

export function useAxios() {
  return client;
}
