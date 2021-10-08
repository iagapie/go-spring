import { SignIn, Tokens } from '@/store/types'
import { requestJSON } from '@/utils/request'
import { endpoints } from '@/utils/constants'

class AuthService {
  public login(dto: SignIn): Promise<Tokens> {
    return requestJSON<any>(endpoints.auth.login, {
      method: 'POST',
      body: JSON.stringify(dto),
    }).then((resp) => ({
      accessToken: resp.body.access_token,
      refreshToken: resp.body.refresh_token,
    }))
  }
}

export default new AuthService()
