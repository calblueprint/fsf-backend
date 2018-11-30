/* originally from crypto/x509/by_file.c */
/* Copyright (C) 1995-1998 Eric Young (eay@cryptsoft.com)
 * All rights reserved.
 *
 * This package is an SSL implementation written
 * by Eric Young (eay@cryptsoft.com).
 * The implementation was written so as to conform with Netscapes SSL.
 *
 * This library is free for commercial and non-commercial use as long as
 * the following conditions are aheared to.  The following conditions
 * apply to all code found in this distribution, be it the RC4, RSA,
 * lhash, DES, etc., code; not just the SSL code.  The SSL documentation
 * included with this distribution is covered by the same copyright terms
 * except that the holder is Tim Hudson (tjh@cryptsoft.com).
 *
 * Copyright remains Eric Young's, and as such any Copyright notices in
 * the code are not to be removed.
 * If this package is used in a product, Eric Young should be given attribution
 * as the author of the parts of the library used.
 * This can be in the form of a textual message at program startup or
 * in documentation (online or textual) provided with the package.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 * 3. All advertising materials mentioning features or use of this software
 *    must display the following acknowledgement:
 *    "This product includes cryptographic software written by
 *     Eric Young (eay@cryptsoft.com)"
 *    The word 'cryptographic' can be left out if the rouines from the library
 *    being used are not cryptographic related :-).
 * 4. If you include any Windows specific code (or a derivative thereof) from
 *    the apps directory (application code) you must include an acknowledgement:
 *    "This product includes software written by Tim Hudson (tjh@cryptsoft.com)"
 *
 * THIS SOFTWARE IS PROVIDED BY ERIC YOUNG ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 *
 * The licence and distribution terms for any publically available version or
 * derivative of this code cannot be changed.  i.e. this code cannot simply be
 * copied and put under another distribution licence
 * [including the GNU Public Licence.]
 */

/* Added code to use a memory object as opposed to a file for lookup purposes. */

#include <string.h>

#include <openssl/bio.h>
#include <openssl/err.h>
#include <openssl/pem.h>
#include <openssl/ssl.h>
#include <openssl/x509_vfy.h>

static int by_mem_ctrl(X509_LOOKUP *, int, const char *, long, char **);
static int X509_load_cert_crl_mem(X509_LOOKUP *ctx, const char *str, int type);

X509_LOOKUP_METHOD x509_mem_lookup =
{
	"Memory from user provided string",
	NULL,
	NULL,
	NULL,
	NULL,
	by_mem_ctrl,
	NULL,
	NULL,
	NULL,
	NULL
};

static X509_LOOKUP_METHOD * X509_LOOKUP_mem()
{
	return (&x509_mem_lookup);
}

static int by_mem_ctrl(X509_LOOKUP *ctx, int cmd, const char * str, long argl, char ** ret)
{
	int status = 0;
	if (cmd == X509_L_FILE_LOAD)
	{
		if (argl == X509_FILETYPE_PEM)
			return X509_load_cert_crl_mem(ctx, str, X509_FILETYPE_PEM) != 0;
	}
	return status;
}

static int X509_load_cert_crl_mem(X509_LOOKUP *ctx, const char *str, int type)
{
	BIO *in = NULL;
	int status = 0;
	int count = 0;
	int i;
	X509 *x = NULL;

	if (str == NULL) return 1;
	in = BIO_new(BIO_s_mem());
	if ((in == NULL) || (BIO_write(in, str, strlen(str)) != strlen(str)) || type != X509_FILETYPE_PEM)
	{
		/* this error isn't the same as the real file one, but it'll serve the same purpose */
		X509err(X509_F_X509_LOAD_CERT_FILE,ERR_R_SYS_LIB);
		goto err;
	}

	for (;;)
	{
		x=PEM_read_bio_X509_AUX(in,NULL,NULL,NULL);
		if (x == NULL)
		{
			if ((ERR_GET_REASON(ERR_peek_last_error()) ==
				PEM_R_NO_START_LINE) && (count > 0))
			{
				ERR_clear_error();
				break;
			}
			else
			{
				X509err(X509_F_X509_LOAD_CERT_FILE,
				ERR_R_PEM_LIB);
				goto err;
			}
		}
		i=X509_STORE_add_cert(ctx->store_ctx,x);
		if (!i) goto err;
		count++;
		X509_free(x);
		x=NULL;
	}
	status=count;

err:
	if (x != NULL) X509_free(x);
	if (in != NULL) BIO_free(in);
	return status;

}

static int X509_STORE_load_mem(X509_STORE *ctx, const char *str)
{
	X509_LOOKUP *lookup;
	if (!str) return 1;
	
	lookup = X509_STORE_add_lookup(ctx, X509_LOOKUP_mem());
	if (lookup == NULL) return 0;
	if (X509_LOOKUP_ctrl(lookup,X509_L_FILE_LOAD,str,X509_FILETYPE_PEM, NULL) != 1)
		return 0;
	return 1;
}
int SSL_CTX_load_verify_locations_mem(SSL_CTX * ctx, const char *str)
{
	return X509_STORE_load_mem(ctx->cert_store, str);
}
