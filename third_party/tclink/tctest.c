/* tctest.c - Simple example of TCLink usage.
 *
 * TCLink Copyright (c) 2013 TrustCommerce.
 * http://www.trustcommerce.com
 * techsupport@trustcommerce.com
 * (949) 387-3747
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 59 Temple Place, Suite 330, Boston, MA  02111-1307  USA
 */

#include "tclink.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>

/* Insert your custid and password below.
 */

#define CUSTID   "TestMerchant"
#define PASSWORD "password"

int main()
{
	char buf[1024];

	TCLinkHandle handle = TCLinkCreate();

	printf("Using TCLink version %s\n", TCLinkGetVersion(buf));

	TCLinkPushParam(handle, "custid",   CUSTID);
	TCLinkPushParam(handle, "password", PASSWORD);
	TCLinkPushParam(handle, "action",   "sale");
	TCLinkPushParam(handle, "amount",   "100");                /* $1.00 */
	TCLinkPushParam(handle, "cc",       "4111111111111111");   /* test Visa card */
	TCLinkPushParam(handle, "exp",      "0404");

	TCLinkSend(handle);

	printf("%s", TCLinkGetEntireResponse(handle, buf, sizeof(buf)));
	
	TCLinkDestroy(handle);

	return 0;
}

