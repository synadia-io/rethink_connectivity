export function parseJWT(token: string): any {
  let base64Url = token.split('.')[1]; // get the payload
  let base64 = base64Url.replace('-', '+').replace('_', '/');
  return JSON.parse(window.atob(base64));
}
