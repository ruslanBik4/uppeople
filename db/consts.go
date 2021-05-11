package db

const sqlGetTypes = "SELECT typname, oid FROM pg_type WHERE typname::text=ANY($1)"

// Table names
const TableTags = "tags"
const TableCandidates = "candidates"
const TableCandidatesOnVacanciesFreelancers = "candidates_on_vacancies_freelancers"
const TableCandidatesToCompanies = "candidates_to_companies"
const TableCandidatesToSkills = "candidates_to_skills"
const TableChromeExtensionSelector = "chrome_extention_selectors"
const TableComments = "comments"
const TableCommentsForCandidates = "comments_for_candidates"
const TableCommentsForCompanies = "comments_for_companies"
const TableCompanies = "companies"
const TableContacts = "contacts"
const TableContactsToPlatforms = "contacts_to_platforms"
const TableDictionary = "dictionary"
const TableEmailTemplates = "email_templates"
const TableFAQCategories = "faq_categories"
const TableFAQQuestions = "faq_qustions"
const TableFavoritesPlatforms = "favorites_platforms"
const TableFirma = "firma"
const TableFormBlocks = "form_blocks"
const TableFormFields = "form_fields"
const TableForms = "forms"
const TableFreelancersVacancies = "freelancers_vacancies"
const TableHomepages = "homepages"
const TableInBoxes = "in_boxes"
const TableIntRevCandidates = "int_rev_candidates"
const TableLanguages = "languages"
const TableLinkedinCandidates = "linkedin_candidates"
const TableLinks = "links"
const TableLocationForVacancies = "location_for_vacancies"
const TableLogs = "logs"
const TableManagement = "management"
const TableMeetings = "meetings"
const TableMigrations = "migrations"
const TableOauthAccessTokens = "oauth_access_tokens"
const TableOauthAuthCodes = "oauth_auth_codes"
const TableOauthClients = "oauth_clients"
const TableOauthPersonalAccessClients = "oauth_personal_access_clients"
const TableOauthRefreshTokens = "oauth_refresh_tokens"
const TablePartners = "partners"
const TablePasswordResets = "password_resets"
const TablePatternsList = "patterns_list"
const TablePlatforms = "platforms"
const TableSentEmails = "sended_emails"
const TableSeniorities = "seniorities"
const TableSkills = "skills"
const TableStatusForVacs = "status_for_vacs"
const TableStatuses = "statuses"
const TableStopLists = "stoplists"
const TableTodoLists = "todolists"
const TableUserToPlatforms = "user_to_platforms"
const TableUserToVacancies = "user_to_vacancies"
const TableUsers = "users"
const TableVacancies = "vacancies"
const TableVacanciesToCandidates = "vacancies_to_candidates"
const TableWPCommentMeta = "wp_commentmeta"
const TableWPComments = "wp_comments"
const TableWPLinks = "wp_links"
const TableWPOptions = "wp_options"
const TableWPPostMeta = "wp_postmeta"
const TableWPPosts = "wp_posts"
const TableWPSimplyStaticPages = "wp_simply_static_pages"
const TableWPTermRelationships = "wp_term_relationships"
const TableWPTermTaxonomy = "wp_term_taxonomy"
const TableWPTermMeta = "wp_termmeta"
const TableWPTerms = "wp_terms"
const TableWPUserMeta = "wp_usermeta"
const TableWPUsers = "wp_users"

//Tag values
const TagFirstContact = "first contact"
const TagInterested = "interested"
const TagReject = "reject"
const TagNoAnswer = "no answer"
const TagClosedToOffers = "closed to offers"
const TagLowSalary = "low salary rate"
const TagWasContactedEarlier = "was contacted earlier"
const TagDoesNotLikeProject = "does not like the project"
const TagTermsDoNotFit = "terms don’t fit"
const TagRemoteOnly = "remote only"
const TagDoesNotFit = "does not fit"

var tagIds TagIdMap
