import { BackendUser } from '@/store/types'
import { requestJSON } from '@/utils/request'
import { endpoints } from '@/utils/constants'

class UserService {
  public me(): Promise<BackendUser> {
    return requestJSON<any>(endpoints.user.me, {
      token: true,
    }).then((resp) => ({
      uuid: resp.body.uuid,
      name: resp.body.name,
      email: resp.body.email,
      createdAt: resp.body.created_at,
      updatedAt: resp.body.updated_at,
    }))
  }
}

export default new UserService()
