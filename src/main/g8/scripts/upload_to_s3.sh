#/bin/bash
service=$name$
file=\${CIRCLE_ARTIFACTS}/snapshot/\${service}_linux_amd64.tar.gz
sha=`echo \${CIRCLE_SHA1} | cut -c1-6`
target=\${service}/\${service}_\${sha}.tar.gz
bucket=$s3_bucket$
resource="/\${bucket}/\${target}"
contentType="application/x-compressed-tar"
dateValue=`date -R`
stringToSign="PUT\n\n\${contentType}\n\${dateValue}\n\${resource}"
s3Key=\$S3_KEY
s3Secret=\$S3_SECRET
signature=`echo -en \${stringToSign} | openssl sha1 -hmac \${s3Secret} -binary | base64`
response=`curl -s -o /dev/null -w "%{http_code}" -X PUT -T "\${file}" \
  -H "Host: \${bucket}.s3.amazonaws.com" \
  -H "Date: \${dateValue}" \
  -H "Content-Type: \${contentType}" \
  -H "Authorization: AWS \${s3Key}:\${signature}" \
  https://\${bucket}.s3.amazonaws.com/\${target}`
echo "Upload status code: \${response}"
if [ \$response -ne "200" ] ; then exit 1 ; fi