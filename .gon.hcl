# gon.hcl
#
# The path follows a pattern
# ./dist/BUILD-ID_TARGET/BINARY-NAME
# source = ["./dist/${artifact}.dmg"]
source = ["./dist/mt-mac_darwin_amd64_v1/malamtime", "./dist/mt-mac_darwin_arm64/malamtime"]
bundle_id = "com.annatarhe.malamtime.cli"

apple_id {
  username = "iamhele1994@gmail.com"
  password = "@env:AC_PASSWORD" # omit this line if env already set
  provider = "VY3AP6C55K"
}

sign {
  application_identity = "Developer ID Application: Le He"
}

zip {
  output_path = "./dist/macos-mt.zip"
}
