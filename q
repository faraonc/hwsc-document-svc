[1mdiff --git a/Gopkg.lock b/Gopkg.lock[m
[1mindex aaf619e..59280d6 100644[m
[1m--- a/Gopkg.lock[m
[1m+++ b/Gopkg.lock[m
[36m@@ -158,45 +158,6 @@[m
   revision = "2a4da1fca4d75c5e1fb254a60b4aa987c73fab08"[m
   version = "v0.14.0"[m
 [m
[31m-[[projects]][m
[31m-  digest = "1:9be9701967c831fb40a4b8b35c474b0a6e3af1f68ac35ad5fde6c44d71cbf4d1"[m
[31m-  name = "github.com/mongodb/mongo-go-driver"[m
[31m-  packages = [[m
[31m-    "bson",[m
[31m-    "bson/bsoncodec",[m
[31m-    "bson/bsonrw",[m
[31m-    "bson/bsontype",[m
[31m-    "bson/primitive",[m
[31m-    "event",[m
[31m-    "internal",[m
[31m-    "mongo",[m
[31m-    "mongo/options",[m
[31m-    "mongo/readconcern",[m
[31m-    "mongo/readpref",[m
[31m-    "mongo/writeconcern",[m
[31m-    "tag",[m
[31m-    "version",[m
[31m-    "x/bsonx",[m
[31m-    "x/bsonx/bsoncore",[m
[31m-    "x/mongo/driver",[m
[31m-    "x/mongo/driver/auth",[m
[31m-    "x/mongo/driver/auth/internal/gssapi",[m
[31m-    "x/mongo/driver/session",[m
[31m-    "x/mongo/driver/topology",[m
[31m-    "x/mongo/driver/uuid",[m
[31m-    "x/network/address",[m
[31m-    "x/network/command",[m
[31m-    "x/network/compressor",[m
[31m-    "x/network/connection",[m
[31m-    "x/network/connstring",[m
[31m-    "x/network/description",[m
[31m-    "x/network/result",[m
[31m-    "x/network/wiremessage",[m
[31m-  ][m
[31m-  pruneopts = "UT"[m
[31m-  revision = "e1f3d104525e687eda76125707c9852e429f2730"[m
[31m-  version = "v0.3.0"[m
[31m-[m
 [[projects]][m
   digest = "1:0028cb19b2e4c3112225cd871870f2d9cf49b9b4276531f03438a88e94be86fe"[m
   name = "github.com/pmezard/go-difflib"[m
[36m@@ -237,6 +198,45 @@[m
   pruneopts = "UT"[m
   revision = "73f8eece6fdcd902c185bf651de50f3828bed5ed"[m
 [m
[32m+[m[32m[[projects]][m
[32m+[m[32m  digest = "1:7280a69811fe89d769f5cbbee7ec10cc33488d1e0d573c6d65e35a8887ab3420"[m
[32m+[m[32m  name = "go.mongodb.org/mongo-driver"[m
[32m+[m[32m  packages = [[m
[32m+[m[32m    "bson",[m
[32m+[m[32m    "bson/bsoncodec",[m
[32m+[m[32m    "bson/bsonrw",[m
[32m+[m[32m    "bson/bsontype",[m
[32m+[m[32m    "bson/primitive",[m
[32m+[m[32m    "event",[m
[32m+[m[32m    "internal",[m
[32m+[m[32m    "mongo",[m
[32m+[m[32m    "mongo/options",[m
[32m+[m[32m    "mongo/readconcern",[m
[32m+[m[32m    "mongo/readpref",[m
[32m+[m[32m    "mongo/writeconcern",[m
[32m+[m[32m    "tag",[m
[32m+[m[32m    "version",[m
[32m+[m[32m    "x/bsonx",[m
[32m+[m[32m    "x/bsonx/bsoncore",[m
[32m+[m[32m    "x/mongo/driver",[m
[32m+[m[32m    "x/mongo/driver/auth",[m
[32m+[m[32m    "x/mongo/driver/auth/internal/gssapi",[m
[32m+[m[32m    "x/mongo/driver/session",[m
[32m+[m[32m    "x/mongo/driver/topology",[m
[32m+[m[32m    "x/mongo/driver/uuid",[m
[32m+[m[32m    "x/network/address",[m
[32m+[m[32m    "x/network/command",[m
[32m+[m[32m    "x/network/compressor",[m
[32m+[m[32m    "x/network/connection",[m
[32m+[m[32m    "x/network/connstring",[m
[32m+[m[32m    "x/network/description",[m
[32m+[m[32m    "x/network/result",[m
[32m+[m[32m    "x/network/wiremessage",[m
[32m+[m[32m  ][m
[32m+[m[32m  pruneopts = "UT"[m
[32m+[m[32m  revision = "ccf36d0607fa2f6f4ae9645e625a7f33b26cf6d1"[m
[32m+[m[32m  version = "v1.0.0-rc1"[m
[32m+[m
 [[projects]][m
   branch = "master"[m
   digest = "1:f92f6956e4059f6a3efc14924d2dd58ba90da25cc57fe07ae3779ef2f5e0c5f2"[m
[36m@@ -368,12 +368,12 @@[m
     "github.com/micro/go-config",[m
     "github.com/micro/go-config/source/env",[m
     "github.com/micro/go-config/source/file",[m
[31m-    "github.com/mongodb/mongo-go-driver/bson",[m
[31m-    "github.com/mongodb/mongo-go-driver/bson/primitive",[m
[31m-    "github.com/mongodb/mongo-go-driver/mongo",[m
[31m-    "github.com/mongodb/mongo-go-driver/mongo/options",[m
     "github.com/segmentio/ksuid",[m
     "github.com/stretchr/testify/assert",[m
[32m+[m[32m    "go.mongodb.org/mongo-driver/bson",[m
[32m+[m[32m    "go.mongodb.org/mongo-driver/bson/primitive",[m
[32m+[m[32m    "go.mongodb.org/mongo-driver/mongo",[m
[32m+[m[32m    "go.mongodb.org/mongo-driver/mongo/options",[m
     "golang.org/x/net/context",[m
     "google.golang.org/grpc",[m
     "google.golang.org/grpc/codes",[m
[1mdiff --git a/Gopkg.toml b/Gopkg.toml[m
[1mindex 197233c..e71ec36 100644[m
[1m--- a/Gopkg.toml[m
[1m+++ b/Gopkg.toml[m
[36m@@ -45,10 +45,6 @@[m
   name = "github.com/micro/go-config"[m
   version = "0.14.0"[m
 [m
[31m-[[constraint]][m
[31m-  name = "github.com/mongodb/mongo-go-driver"[m
[31m-  version = "0.3.0"[m
[31m-[m
 [[constraint]][m
   name = "github.com/segmentio/ksuid"[m
   version = "1.0.2"[m
[36m@@ -65,6 +61,10 @@[m
   name = "google.golang.org/grpc"[m
   version = "1.18.0"[m
 [m
[32m+[m[32m[[constraint]][m
[32m+[m[32m  name = "go.mongodb.org/mongo-driver"[m
[32m+[m[32m  version = "1.0.0-rc1"[m
[32m+[m
 [prune][m
   go-tests = true[m
   unused-packages = true[m
[1mdiff --git a/README.md b/README.md[m
[1mindex 0931967..4bc6783 100644[m
[1m--- a/README.md[m
[1m+++ b/README.md[m
[36m@@ -36,7 +36,7 @@[m [mThe proto file and compiled proto buffers are located in [hwsc-api-blocks](https[m
 - GoLang version [go 1.11.5](https://golang.org/dl/)[m
 - GoLang Dependency Management [dep](https://github.com/golang/dep)[m
 - Go Source Code Linter [golint](https://github.com/golang/lint)[m
[31m-- mongo-go-driver beta [0.3.0](https://github.com/mongodb/mongo-go-driver)[m
[32m+[m[32m- mongo-go-driver beta [1.0.0](https://github.com/mongodb/mongo-go-driver)[m
 - Docker[m
 - [Optional] If a new proto file and compiled proto buffer exists in [hwsc-api-blocks](https://github.com/hwsc-org/hwsc-api-blocks/tree/master/int/hwsc-document-svc/proto), update dependency ``$dep ensure -update``[m
 [m
[1mdiff --git a/service/db.go b/service/db.go[m
[1mindex a268204..22d8527 100644[m
[1m--- a/service/db.go[m
[1m+++ b/service/db.go[m
[36m@@ -4,7 +4,8 @@[m [mimport ([m
 	"github.com/hwsc-org/hwsc-document-svc/conf"[m
 	"github.com/hwsc-org/hwsc-document-svc/consts"[m
 	log "github.com/hwsc-org/hwsc-lib/logger"[m
[31m-	"github.com/mongodb/mongo-go-driver/mongo"[m
[32m+[m	[32m"go.mongodb.org/mongo-driver/mongo"[m
[32m+[m	[32m"go.mongodb.org/mongo-driver/mongo/options"[m
 	"golang.org/x/net/context"[m
 	"os"[m
 	"os/signal"[m
[36m@@ -41,7 +42,7 @@[m [mfunc init() {[m
 // dialMongoDB connects a client to MongoDB server.[m
 // Returns a MongoDB Client or any dialing error.[m
 func dialMongoDB(uri *string) (*mongo.Client, error) {[m
[31m-	client, err := mongo.Connect(context.TODO(), *uri)[m
[32m+[m	[32mclient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(*uri))[m
 	if err != nil {[m
 		return nil, err[m
 	}[m
[1mdiff --git a/service/db_test.go b/service/db_test.go[m
[1mindex 4c6bf08..44293dd 100644[m
[1m--- a/service/db_test.go[m
[1m+++ b/service/db_test.go[m
[36m@@ -14,7 +14,10 @@[m [mfunc TestDialMongoDB(t *testing.T) {[m
 		errorStr string[m
 	}{[m
 		{conf.DocumentDB.Reader, false, ""},[m
[31m-		{"", true, "error parsing uri (): scheme must be \"mongodb\" or \"mongodb+srv\""},[m
[32m+[m		[32m{"", true, "server selection error: server selection timeout\ncurrent topology: " +[m
[32m+[m			[32m"Type: Unknown\nServers:\nAddr: localhost:27017, Type: Unknown, State: Connected, Avergage RTT: 0, " +[m
[32m+[m			[32m"Last error: dial tcp [::1]:27017: connectex: No connection could be made because the target machine " +[m
[32m+[m			[32m"actively refused it.\n"},[m
 	}[m
 [m
 	for _, c := range cases {[m
[1mdiff --git a/service/service.go b/service/service.go[m
[1mindex c96ede9..110f977 100644[m
[1m--- a/service/service.go[m
[1m+++ b/service/service.go[m
[36m@@ -9,8 +9,8 @@[m [mimport ([m
 	"github.com/hwsc-org/hwsc-document-svc/consts"[m
 	log "github.com/hwsc-org/hwsc-lib/logger"[m
 	"github.com/kylelemons/godebug/pretty"[m
[31m-	"github.com/mongodb/mongo-go-driver/bson"[m
[31m-	"github.com/mongodb/mongo-go-driver/mongo/options"[m
[32m+[m	[32m"go.mongodb.org/mongo-driver/bson"[m
[32m+[m	[32m"go.mongodb.org/mongo-driver/mongo/options"[m
 	"golang.org/x/net/context"[m
 	"google.golang.org/grpc/codes"[m
 	"google.golang.org/grpc/status"[m
[1mdiff --git a/service/utility.go b/service/utility.go[m
[1mindex 4f7e773..6fffc2b 100644[m
[1m--- a/service/utility.go[m
[1m+++ b/service/utility.go[m
[36m@@ -6,9 +6,9 @@[m [mimport ([m
 	pbdoc "github.com/hwsc-org/hwsc-api-blocks/lib"[m
 	"github.com/hwsc-org/hwsc-document-svc/consts"[m
 	log "github.com/hwsc-org/hwsc-lib/logger"[m
[31m-	"github.com/mongodb/mongo-go-driver/bson"[m
[31m-	"github.com/mongodb/mongo-go-driver/bson/primitive"[m
 	"github.com/segmentio/ksuid"[m
[32m+[m	[32m"go.mongodb.org/mongo-driver/bson"[m
[32m+[m	[32m"go.mongodb.org/mongo-driver/bson/primitive"[m
 	"net/http"[m
 	"net/url"[m
 	"regexp"[m
[1mdiff --git a/service/utility_test.go b/service/utility_test.go[m
[1mindex b73ff44..1e9416a 100644[m
[1m--- a/service/utility_test.go[m
[1m+++ b/service/utility_test.go[m
[36m@@ -3,8 +3,8 @@[m [mpackage service[m
 import ([m
 	pbdoc "github.com/hwsc-org/hwsc-api-blocks/lib"[m
 	"github.com/hwsc-org/hwsc-document-svc/consts"[m
[31m-	"github.com/mongodb/mongo-go-driver/bson"[m
 	"github.com/stretchr/testify/assert"[m
[32m+[m	[32m"go.mongodb.org/mongo-driver/bson"[m
 	"sync"[m
 	"testing"[m
 	"time"[m
