import { library } from '@fortawesome/fontawesome-svg-core'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { faInstagram, faTwitter } from '@fortawesome/free-brands-svg-icons'
import FontAwesomeIcon from '@/components/shared/FontAwesomeIcon.vue'

library.add(fas, faTwitter, faInstagram)

export { FontAwesomeIcon }
