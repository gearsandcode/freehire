package skilltag

// wordAliases maps a bare alphanumeric token to its canonical skill slug. Matched
// on whole word tokens (the word pass), so an entry never matches inside a larger
// word. Ambiguous English words (go, c, r) are deliberately ABSENT here — they are
// emitted only via an unambiguous alias (golang) or a phrase ("c developer"), so
// "please go home" or "plan c" never tags. Seed list; expand toward ~200 from
// MIND-tech-ontology and github/linguist languages.yml — pure data, no engine change.
var wordAliases = map[string]string{
	// languages (unambiguous)
	"golang": "go", "python": "python", "java": "java",
	"javascript": "javascript",
	"typescript": "typescript", "ts": "typescript",
	"rust": "rust", "kotlin": "kotlin", "swift": "swift", "scala": "scala",
	"ruby": "ruby", "php": "php", "elixir": "elixir", "erlang": "erlang",
	"clojure": "clojure", "haskell": "haskell", "perl": "perl", "lua": "lua",
	"dart": "dart", "groovy": "groovy", "solidity": "solidity",
	// frontend
	"react": "react", "reactjs": "react",
	"angular": "angular", "angularjs": "angular",
	"vue": "vue", "vuejs": "vue", "svelte": "svelte", "nextjs": "nextjs", "nodejs": "nodejs",
	"nuxt": "nuxt", "redux": "redux", "tailwind": "tailwind",
	"webpack": "webpack", "vite": "vite",
	// backend frameworks
	"django": "django", "flask": "flask", "fastapi": "fastapi",
	"spring": "spring", "rails": "rails", "laravel": "laravel",
	"symfony": "symfony", "express": "express", "nestjs": "nestjs",
	"gin": "gin", "fiber": "fiber", "dotnet": "dotnet",
	// datastores
	"postgres": "postgresql", "postgresql": "postgresql", "psql": "postgresql",
	"mysql": "mysql", "mariadb": "mariadb", "sqlite": "sqlite",
	"redis": "redis", "memcached": "memcached",
	"mongodb": "mongodb", "mongo": "mongodb", "cassandra": "cassandra",
	"dynamodb": "dynamodb", "elasticsearch": "elasticsearch",
	"clickhouse": "clickhouse", "snowflake": "snowflake",
	"kafka": "kafka", "rabbitmq": "rabbitmq", "nats": "nats",
	// infra / cloud
	"kubernetes": "kubernetes", "k8s": "kubernetes", "docker": "docker",
	"terraform": "terraform", "ansible": "ansible", "pulumi": "pulumi",
	"helm": "helm", "aws": "aws", "gcp": "gcp", "azure": "azure",
	"linux": "linux", "nginx": "nginx",
	"prometheus": "prometheus", "grafana": "grafana", "jenkins": "jenkins",
	// api / data
	"graphql": "graphql", "grpc": "grpc", "rest": "rest",
	"pytorch": "pytorch", "tensorflow": "tensorflow", "pandas": "pandas",
	"numpy": "numpy", "spark": "spark", "airflow": "airflow", "dbt": "dbt",
}

// phraseAlias is a punctuated or multi-word term matched against the normalized
// text with whole-term boundaries (the phrase pass). Canonicals are facet-safe
// slugs (cpp, csharp, ci-cd) rather than the raw punctuated form.
type phraseAlias struct {
	alias     string
	canonical string
}

// phraseAliases covers terms the word pass cannot see because they contain
// non-alphanumeric characters or spaces. Includes the ONLY routes by which an
// ambiguous canonical (c) may be emitted.
var phraseAliases = []phraseAlias{
	{"c++", "cpp"}, {"c/c++", "cpp"},
	{"c#", "csharp"},
	{".net", "dotnet"}, {"asp.net", "dotnet"},
	{"node.js", "nodejs"}, {"node js", "nodejs"},
	{"next.js", "nextjs"},
	{"vue.js", "vue"},
	{"react native", "react-native"},
	{"react.js", "react"},
	{"objective-c", "objective-c"},
	{"ci/cd", "ci-cd"}, {"ci cd", "ci-cd"},
	{"c developer", "c"}, {"c programming", "c"}, {"ansi c", "c"},
	{"machine learning", "machine-learning"},
}
