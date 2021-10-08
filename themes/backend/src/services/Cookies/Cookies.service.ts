import Cookie from 'js-cookie'

Cookie.defaults = {
  sameSite: 'Lax',
}

class CookiesService {
  public static get<T>(key: string, defaultValue: T): T {
    const item = Cookie.get(key)

    return item ? JSON.parse(item) : defaultValue
  }

  public static set(key: string, value: any): void {
    Cookie.set(key, JSON.stringify(value))
  }

  public static remove(key: string): void {
    Cookie.remove(key)
  }
}

export default CookiesService
