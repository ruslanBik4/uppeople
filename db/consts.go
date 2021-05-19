package db

const SQL_GET_TYPES = "SELECT typname, oid FROM pg_type WHERE typname::text=ANY($1)"

// Table names
const TABLE_TAGS = "tags"
const TABLE_CANDIDATES = "candidates"
const TABLE_CANDIDATES_ON_VACANCIES_FREELANCERS = "candidates_on_vacancies_freelancers"
const TABLE_CANDIDATES_TO_COMPANIES = "candidates_to_companies"
const TABLE_CANDIDATES_TO_SKILLS = "candidates_to_skills"
const TABLE_CHROME_EXTENTION_SELECTORS = "chrome_extention_selectors"
const TABLE_COMMENTS = "comments"
const TABLE_COMMENTS_FOR_CANDIDATES = "comments_for_candidates"
const TABLE_COMMENTS_FOR_COMPANIES = "comments_for_companies"
const TABLE_COMPANIES = "companies"
const TABLE_CONTACTS = "contacts"
const TABLE_CONTACTS_TO_PLATFORMS = "contacts_to_platforms"
const TABLE_DICTIONARY = "dictionary"
const TABLE_EMAIL_TEMPLATES = "email_templates"
const TABLE_FAQ_CATEGORIES = "faq_categories"
const TABLE_FAQ_QUESTIONS = "faq_qustions"
const TABLE_FAVORITE_PLATFORMS = "favorites_platforms"
const TABLE_FIRMA = "firma"
const TABLE_FORM_BLOCKS = "form_blocks"
const TABLE_FORM_FIELDS = "form_fields"
const TABLE_FORMS = "forms"
const TABLE_FREELANCER_VACANCIES = "freelancers_vacancies"
const TABLE_HOMEPAGES = "homepages"
const TABLE_IN_BOXES = "in_boxes"
const TABLE_INT_REV_CANDIDATES = "int_rev_candidates"
const TABLE_LANGUAGES = "languages"
const TABLE_LINKEDIN_CANDIDATES = "linkedin_candidates"
const TABLE_LINKS = "links"
const TABLE_LOCATION_FOR_VACANCIES = "location_for_vacancies"
const TABLE_LOGS = "logs"
const TABLE_MANAGEMENT = "management"
const TABLE_MEETINGS = "meetings"
const TABLE_MIGRATIONS = "migrations"
const TABLE_OAUTH_ACCESS_TOKENS = "oauth_access_tokens"
const TABLE_OAUTH_AUTH_CODES = "oauth_auth_codes"
const TABLE_OAUTH_CLIENTS = "oauth_clients"
const TABLE_OAUTH_PERSONAL_ACCESS_CLIENTS = "oauth_personal_access_clients"
const TABLE_OAUTH_REFRESH_TOKENS = "oauth_refresh_tokens"
const TABLE_PARTNERS = "partners"
const TABLE_PASSWORD_RESETS = "password_resets"
const TABLE_PATTERN_LIST = "patterns_list"
const TABLE_PLATFORMS = "platforms"
const TABLE_SENT_EMAILS = "sended_emails"
const TABLE_SENIORITIES = "seniorities"
const TABLE_SKILLS = "skills"
const TABLE_STATUS_FOR_VACS = "status_for_vacs"
const TABLE_STATUSES = "statuses"
const TABLE_STOP_LISTS = "stoplists"
const TABLE_TODO_LISTS = "todolists"
const TABLE_USER_TO_PLATFORMS = "user_to_platforms"
const TABLE_USER_TO_VACANCIES = "user_to_vacancies"
const TABLE_USERS = "users"
const TABLE_VACANCIES = "vacancies"
const TABLE_VACANCIES_TO_CANDIDATES = "vacancies_to_candidates"
const TABLE_WP_COMMENT_META = "wp_commentmeta"
const TABLE_WP_COMMENTS = "wp_comments"
const TABLE_WP_LINKS = "wp_links"
const TABLE_WP_OPTIONS = "wp_options"
const TABLE_WP_POST_META = "wp_postmeta"
const TABLE_WP_POSTS = "wp_posts"
const TABLE_WP_SIMPLY_STATIC_PAGES = "wp_simply_static_pages"
const TABLE_WP_TERM_RELATIONSHIPS = "wp_term_relationships"
const TABLE_WP_TERM_TAXONOMY = "wp_term_taxonomy"
const TABLE_WP_TERM_META = "wp_termmeta"
const TABLE_WP_TERMS = "wp_terms"
const TABLE_WP_USER_META = "wp_usermeta"
const TABLE_WP_USERS = "wp_users"

//Tag values
const TAG_FIRST_CONTACT = "first contact"
const TAG_INTERESTED = "interested"
const TAG_REJECT = "reject"
const TAG_NO_ANSWER = "no answer"
const TAG_CLOSED_TO_OFFERS = "closed to offers"
const TAG_LOW_SALARY = "low salary rate"
const TAG_WAS_CONTACTED_EARLIER = "was contacted earlier"
const TAG_DOES_NOT_LIKE_PROJECT = "does not like the project"
const TAG_TERMS_DO_NOT_FIT = "terms don’t fit"
const TAG_REMOTE_ONLY = "remote only"
const TAG_DOES_NOT_FIT = "does not fit"

//Status values
const STATUS_HOT = "Hot"
const STATUS_OPEN = "Open"
const STATUS_CLOSED = "Closed"
const STATUS_PAUSED = "Paused"

//Status for vacancies values
const STATUS_FOR_VAC_INTERVIEW = "Interview"
const STATUS_FOR_VAC_TEST = "Test"
const STATUS_FOR_VAC_FINAL_INTERVIEW = "Final Interview"
const STATUS_FOR_VAC_OFFER = "OFFER"
const STATUS_FOR_VAC_HIRED = "Hired"
const STATUS_FOR_VAC_WR = "WR"
const STATUS_FOR_VAC_REVIEW = "Review"
const STATUS_FOR_VAC_REJECTED = "Rejected"
const STATUS_FOR_VAC_ON_HOLD = "On hold"

//Seniorities values
const SENIORITY_JUN = "Jun"
const SENIORITY_MID = "Mid"
const SENIORITY_SEN = "Sen"
const SENIORITY_LEAD = "Lead"
const SENIORITY_ARCHITECT = "Architect"
const SENIORITY_JUN_MID = "Jun-Mid"
const SENIORITY_MID_SEN = "Mid-Sen"
const SENIORITY_SEN_LEAD = "Sen-Lead"

//Consts for table values
var (
	tagIds            TagIdMap
	statusesIds       StatusIdMap
	statusesForVacIds StatusForVacIdMap
	seniorityIds      SeniorityIdMap
	platformIds       PlatformsIdMap
)

// Consts forgetter handlerds
var (
	tagIdsAsSU            SelectedUnits
	reasonsIdsAsSU        SelectedUnits
	statusesIdsAsSU       SelectedUnits
	statusesForVacIdsAsSU SelectedUnits
	seniorityIdsAsSU      SelectedUnits
	platformIdsAsSU       SelectedUnits
)
