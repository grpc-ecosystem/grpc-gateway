FROM gitpod/workspace-full

# Install swagger-codegen
ENV SWAGGER_CODEGEN_VERSION=2.4.8
RUN wget https://repo1.maven.org/maven2/io/swagger/swagger-codegen-cli/${SWAGGER_CODEGEN_VERSION}/swagger-codegen-cli-${SWAGGER_CODEGEN_VERSION}.jar \
    -O /home/gitpod/swagger-codegen-cli.jar && \
    echo -e '#!/bin/bash\njava -jar /home/gitpod/swagger-codegen-cli.jar "$@"' > /home/gitpod/swagger-codegen && \
    sudo mv /home/gitpod/swagger-codegen /usr/local/bin/swagger-codegen && \
	sudo chmod +x /usr/local/bin/swagger-codegen
