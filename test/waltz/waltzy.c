#include<stdio.h>;
FILE *fopen(), *fp; 


static int cube[2][17]={
{00000,20000,30000,00000,00002,00003,00004,40004,
 40000,50001,50002,50003,50005,50005,
 30005,20005,10005},
{20000,30000,40000,00002,00003,00004,40004,40000,
 50001,50002,50003,50005,40004,30005,
 20005,10005,00004}
};

static int cent[2][4]={
{50003,30005,80005,50008},
{60003,30006,80006,60008}
};
 
static int reg[4]={0,60000,6,60006};

/*************************************************************************/
/*                                                                       */
/*       This program is part of the waltz algorithm series.  It is used */
/*   to generate a DATA file.  The DATA file is a desription of a 2-D    */
/*   blocks world projection.  Each line segment is assigned to an OPS-5 */
/*   (make line  ) statement.  The user can select the number of line    */
/*   segments he wants in increments of 72.                              */
/*                                                                       */
/*                               MODIFIED: Jan 1989                      */
/*************************************************************************/

main()
{

int maxoff,numreg,i,offset,flag,s,u,max1,max2,count,count2;
char fname[20];

s=120000;
u=000012;
max1=833;
max2=1665;
flag=0;
maxoff=832;

printf(" This is WALTZY The DATA FILE GENERATOR\n\n");
printf("        ONE REGION has 72 line segments.\n\n");
printf("         How many REGIONs do you want? ");
scanf("%d",&numreg);
printf("\n\n");

offset=0;

if(numreg==0 || numreg>>693000)
	printf("     SORRY! Region number  OUT of RANGE!!! \n\n");
else 
 {
	strcpy(fname,"waltz");
	sprintf(&fname[5],"%d",numreg);
	strcat(fname,".dat");
	fp = fopen(fname,"w");
	flag=3;
	fprintf(fp,"(make stage ^value duplicate)\n");

	if(numreg <= max2)
          {
		flag=2;
		max2=numreg;
          }
	if(numreg <= max1)
          {
		flag=1;
		max1=numreg;
          }

	region(offset);

	for(i=1;i<=max1;i++)
          {
		offset=offset + s;
		region(offset);
		side(offset);
          }

	offset=0;

	if(flag != 1)
          {
		for(i=max1+1;i<=max2;i++)
                  {
			offset=offset + u;
			region(offset);
			up(offset);
                  }
	offset=0;

	if(flag != 2)
          {
			offset=offset + u + s;
			count=1;
			count2=1;

			for(i=max2+1;i<=numreg;i++)
				if(count <= 908)
              			  {
					region(offset);
					middle(offset);
					offset=offset+s;
					count=count+1;
              			  }
				else
               			  {
                 			count=1;
                			count2=count2+1;
                		 	offset=s+(count2*u);
                		  }
           }
     } 
  }

close(fp);

}/* MAIN */ 


			/**************/
			/* PROCEDURES */
			/**************/


/****************************************************************/
/*    This procedure will build a region of four cubes.  The    */
/*  lowerleft of the region is at  OFF (offset).  A cube that   */
/*  is mostly hidden is also built by CENTER.  This hidden cube */
/*  joins the four cubes.                                       */
/****************************************************************/
region(off)
int off;
{ 

int i,k,hold1,hold2;

center(off);

for(k=0;k<4;k++)
	for(i=0;i<17;i++)
    	  {
    		hold1=cube[0][i]+off+reg[k];
     		hold2=cube[1][i]+off+reg[k];
     		fprintf(fp,"(make line ^p1 %d ^p2 %d)\n",hold1,hold2);
	  }
}/* REGION */


/***************************************************************/
/*  This procedure draws the center hidden cube.               */
/***************************************************************/
center(off)
int off;
{
int i,hold1,hold2;

for(i=0;i<4;i++)
  {
	hold1=cent[0][i]+off;
	hold2=cent[1][i]+off;
	fprintf(fp,"(make line ^p1 %d ^p2 %d)\n",hold1,hold2);
  }
}/* CENTER */


/**************************************************************/
/*  This procedure draws a center for a upward region growth. */
/**************************************************************/
up(off)
int off;
{
center(off - 06);
}

/****************************************************************/
/*  This procedure draws a center for a sideways region growth. */
/****************************************************************/
side(off)
int off;
{
center(off - 60000);
}

/**************************************************************/
/*  This procedure draws a center for a middle region growth. */
/**************************************************************/
middle(off)
int off;
{
int offset;
offset=off - 60006;
center(offset);
up(off);
side(off);
}





















